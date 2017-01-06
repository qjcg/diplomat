package gopdf

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"strings"
)

func buildImgProp(imginfo imgInfo) (*bytes.Buffer, error) {

	var buffer bytes.Buffer
	buffer.WriteString("<</Type /XObject\n")
	buffer.WriteString("/Subtype /Image\n")
	buffer.WriteString(fmt.Sprintf("/Width %d\n", imginfo.w))  // /Width 675\n"
	buffer.WriteString(fmt.Sprintf("/Height %d\n", imginfo.h)) //  /Height 942\n"
	if isColspaceIndexed(imginfo) {
		size := len(imginfo.pal)/3 - 1
		buffer.WriteString(fmt.Sprintf("/ColorSpace [/Indexed /DeviceRGB %d %d 0 R]\n", size, imginfo.deviceRGBObjID+1))
	} else {
		buffer.WriteString(fmt.Sprintf("/ColorSpace /%s\n", imginfo.colspace))
		if imginfo.colspace == "DeviceCMYK" {
			buffer.WriteString("/Decode [1 0 1 0 1 0 1 0]\n")
		}
	}
	buffer.WriteString(fmt.Sprintf("/BitsPerComponent %s\n", imginfo.bitsPerComponent))
	if strings.TrimSpace(imginfo.filter) != "" {
		buffer.WriteString(fmt.Sprintf("/Filter /%s\n", imginfo.filter))
	}

	if strings.TrimSpace(imginfo.decodeParms) != "" {
		buffer.WriteString(fmt.Sprintf("/DecodeParms <<%s>>\n", imginfo.decodeParms))
	}

	if imginfo.trns != nil && len(imginfo.trns) > 0 {
		j := 0
		max := len(imginfo.trns)
		var trns bytes.Buffer
		for j < max {
			//trn := int(imginfo.trns[j])
			//trns.WriteByte(imginfo.trns[j])
			trns.WriteString(fmt.Sprintf("%d", imginfo.trns[j]))
			trns.WriteString(" ")
			trns.WriteString(fmt.Sprintf("%d", imginfo.trns[j]))
			trns.WriteString(" ")
			j++
		}
		buffer.WriteString(fmt.Sprintf("/Mask [%s]\n", trns.String()))
	}

	if haveSMask(imginfo) {
		buffer.WriteString(fmt.Sprintf("/SMask %d 0 R\n", imginfo.smarkObjID+1))
	}

	return &buffer, nil
}

func isColspaceIndexed(imginfo imgInfo) bool {
	if imginfo.colspace == "Indexed" {
		return true
	}
	return false
}

func haveSMask(imginfo imgInfo) bool {
	if imginfo.smask != nil && len(imginfo.smask) > 0 {
		return true
	}
	return false
}

func parseImg(raw []byte) (imgInfo, error) {
	//fmt.Printf("----------\n")
	var info imgInfo
	imgConfig, formatname, err := image.DecodeConfig(bytes.NewBuffer(raw))
	if err != nil {
		return info, err
	}
	info.formatName = formatname

	if formatname == "jpeg" {

		err = parseImgJpg(&info, imgConfig)
		if err != nil {
			return info, err
		}
		info.data = raw
	} else if formatname == "png" {
		err = paesePng(raw, &info, imgConfig)
		if err != nil {
			return info, err
		}
	}

	//fmt.Printf("%#v\n", info)

	return info, nil
}

func parseImgJpg(info *imgInfo, imgConfig image.Config) error {
	if imgConfig.ColorModel == color.YCbCrModel {
		info.colspace = "DeviceRGB"
	} else if imgConfig.ColorModel == color.GrayModel {
		info.colspace = "DeviceGray"
	} else if imgConfig.ColorModel == color.CMYKModel {
		info.colspace = "DeviceCMYK"
	} else {
		return errors.New("color model not support")
	}
	info.bitsPerComponent = "8"
	info.filter = "DCTDecode"

	info.h = imgConfig.Height
	info.w = imgConfig.Width

	return nil
}

var pngMagicNumber = []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}
var pngIHDR = []byte{0x49, 0x48, 0x44, 0x52}

func paesePng(raw []byte, info *imgInfo, imgConfig image.Config) error {
	f := bytes.NewReader(raw)
	f.Seek(0, 0)
	b, err := readBytes(f, 8)
	if err != nil {
		return err
	}
	if !compareBytes(b, pngMagicNumber) {
		return errors.New("Not a PNG file")
	}

	f.Seek(4, 1) //skip header chunk
	b, err = readBytes(f, 4)
	if err != nil {
		return err
	}
	if !compareBytes(b, pngIHDR) {
		return errors.New("Incorrect PNG file")
	}

	w, err := readInt(f)
	if err != nil {
		return err
	}
	h, err := readInt(f)
	if err != nil {
		return err
	}
	//fmt.Printf("w=%d h=%d\n", w, h)

	bpc, err := readBytes(f, 1)
	if err != nil {
		return err
	}

	if bpc[0] > 8 {
		return errors.New("16-bit depth not supported")
	}

	ct, err := readBytes(f, 1)
	if err != nil {
		return err
	}

	var colspace string
	if ct[0] == 0 || ct[0] == 4 {
		colspace = "DeviceGray"
	} else if ct[0] == 2 || ct[0] == 6 {
		colspace = "DeviceRGB"
	} else if ct[0] == 3 {
		colspace = "Indexed"
	} else {
		return errors.New("Unknown color type")
	}

	compressionMethod, err := readBytes(f, 1)
	if err != nil {
		return err
	}
	if compressionMethod[0] != 0 {
		return errors.New("Unknown compression method")
	}

	filterMethod, err := readBytes(f, 1)
	if err != nil {
		return err
	}
	if filterMethod[0] != 0 {
		return errors.New("Unknown filter method")
	}

	interlacing, err := readBytes(f, 1)
	if err != nil {
		return err
	}
	if interlacing[0] != 0 {
		return errors.New("Interlacing not supported")
	}

	_, err = f.Seek(4, 1) //skip
	if err != nil {
		return err
	}

	//decodeParms := "/Predictor 15 /Colors '.($colspace=='DeviceRGB' ? 3 : 1).' /BitsPerComponent '.$bpc.' /Columns '.$w;

	var pal []byte
	var trns []byte
	var data []byte
	for {
		n, err := readInt(f)
		if err != nil {
			return err
		}

		typ, err := readBytes(f, 4)
		//fmt.Printf(">>>>>%s\n", string(typ))
		if err != nil {
			return err
		}

		if string(typ) == "PLTE" {
			pal, err = readBytes(f, n)
			if err != nil {
				return err
			}
			_, err = f.Seek(int64(4), 1) //skip
			if err != nil {
				return err
			}
		} else if string(typ) == "tRNS" {

			var t []byte
			t, err = readBytes(f, n)
			if err != nil {
				return err
			}

			if ct[0] == 0 {
				trns = []byte{(t[1])}
			} else if ct[0] == 2 {
				trns = []byte{t[1], t[3], t[5]}
			} else {
				pos := strings.Index(string(t), "\x00")
				if pos >= 0 {
					trns = []byte{byte(pos)}
				}
			}

			_, err = f.Seek(int64(4), 1) //skip
			if err != nil {
				return err
			}

		} else if string(typ) == "IDAT" {
			//fmt.Printf("n=%d\n\n", n)
			var d []byte
			d, err = readBytes(f, n)
			if err != nil {
				return err
			}
			data = append(data, d...)
			_, err = f.Seek(int64(4), 1) //skip
			if err != nil {
				return err
			}
		} else if string(typ) == "IEND" {
			break
		} else {
			_, err = f.Seek(int64(n+4), 1) //skip
			if err != nil {
				return err
			}
		}

		if n <= 0 {
			break
		}
	} //end for

	//info.data = data //ok
	info.trns = trns
	info.pal = pal

	//fmt.Printf("data= %x", md5.Sum(data))

	if colspace == "Indexed" && strings.TrimSpace(string(pal)) == "" {
		return errors.New("Missing palette")
	}

	info.w = w
	info.h = h
	info.colspace = colspace
	info.bitsPerComponent = fmt.Sprintf("%d", int(bpc[0]))
	info.filter = "FlateDecode"

	colors := 1
	if colspace == "DeviceRGB" {
		colors = 3
	}
	info.decodeParms = fmt.Sprintf("/Predictor 15 /Colors  %d /BitsPerComponent %s /Columns %d", colors, info.bitsPerComponent, w)

	//fmt.Printf("%d = ct[0]\n", ct[0])
	//fmt.Printf("%x\n", md5.Sum(data))
	if ct[0] >= 4 {
		zipReader, err := zlib.NewReader(bytes.NewReader(data))
		if err != nil {
			return err
		}
		defer zipReader.Close()
		afterZipData, err := ioutil.ReadAll(zipReader)
		if err != nil {
			return err
		}

		var color []byte
		var alpha []byte
		if ct[0] == 4 {
			// Gray image
			length := 2 * w
			i := 0
			for i < h {
				pos := (1 + length) * i
				color = append(color, afterZipData[pos])
				alpha = append(alpha, afterZipData[pos])
				line := afterZipData[pos+1 : pos+length+1]
				j := 0
				max := len(line)
				for j < max {
					color = append(color, line[j])
					j++
					alpha = append(alpha, line[j])
					j++
				}
				i++
			}
			//fmt.Print("aaaaa")

		} else {
			// RGB image
			length := 4 * w
			i := 0
			for i < h {
				pos := (1 + length) * i
				color = append(color, afterZipData[pos])
				alpha = append(alpha, afterZipData[pos])
				line := afterZipData[pos+1 : pos+length+1]
				j := 0
				max := len(line)
				for j < max {
					color = append(color, line[j:j+3]...)
					alpha = append(alpha, line[j+3])
					j = j + 4
				}

				i++
			}
			info.smask, err = compress(alpha)
			if err != nil {
				return err
			}

			info.data, err = compress(color)
			if err != nil {
				return err
			}
		}

	} else {
		info.data = data
	}

	return nil
}

func compress(data []byte) ([]byte, error) {
	var results []byte
	var buff bytes.Buffer
	zwr, err := zlib.NewWriterLevel(&buff, zlib.BestSpeed)
	if err != nil {
		return results, err
	}
	_, err = zwr.Write(data)
	if err != nil {
		return results, err
	}
	zwr.Close()
	return buff.Bytes(), nil
}

func readUInt(f *bytes.Reader) (uint, error) {
	buff, err := readBytes(f, 4)
	//fmt.Printf("%#v\n\n", buff)
	if err != nil {
		return 0, err
	}
	n := binary.BigEndian.Uint32(buff)
	return uint(n), nil
}

func readInt(f *bytes.Reader) (int, error) {

	u, err := readUInt(f)
	if err != nil {
		return 0, err
	}
	var v int
	if u >= 0x8000 {
		v = int(u) - 65536
	} else {
		v = int(u)
	}
	return v, nil
}

func readBytes(f *bytes.Reader, len int) ([]byte, error) {
	b := make([]byte, len)
	_, err := f.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func compareBytes(a []byte, b []byte) bool {
	if a == nil && b == nil {
		return true
	} else if a == nil {
		return false
	} else if b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	i := 0
	max := len(a)
	for i < max {
		if a[i] != b[i] {
			return false
		}
		i++
	}
	return true
}

func isDeviceRGB(formatname string, img *image.Image) bool {
	if _, ok := (*img).(*image.YCbCr); ok {
		return true
	} else if _, ok := (*img).(*image.NRGBA); ok {
		return true
	}
	return false
}

//ImgReactagleToWH  Rectangle to W and H
func ImgReactagleToWH(imageRect image.Rectangle) (float64, float64) {
	k := 1
	w := -128 //init
	h := -128 //init
	if w < 0 {
		w = -imageRect.Dx() * 72 / w / k
	}
	if h < 0 {
		h = -imageRect.Dy() * 72 / h / k
	}
	if w == 0 {
		w = h * imageRect.Dx() / imageRect.Dy()
	}
	if h == 0 {
		h = w * imageRect.Dy() / imageRect.Dx()
	}
	return float64(w), float64(h)
}
