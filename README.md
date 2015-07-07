Generate PDF diplomas.

# Features
- takes as argument:
    - instructor
    - course title
    - start / end dates
    - student names

# Usage

Passing arguments via command line flags:

```
godip -i 'Jerry Grapes' -c 'Doing groovy stuff with TempleOS' -d 'May 5-7 2015 (22.5h)' -s 'Jack Groover','Jerry Funkytown', 'Janine Frankfurter'
```

Reading parameters from config file:

```
godip -f config.json
```
