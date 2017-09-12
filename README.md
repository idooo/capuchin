# Capuchin

Small and simple Chaos Monkey inspired by Netflix Simian Army

- Can terminate tagged instances in autoscaling groups
- Can stop tagged instances
- Can restore stopped instances
- Sends notifications to Cloudwatch Logs Streams (capuchin-log-group -> capuchin-log-stream), creates it if needed

## install

You have to install one dependency:

```
go get -u github.com/aws/aws-sdk-go
```

## Configuration 

Example configuration file with specified tags:

```
{
    "Restore": {
        "Instances": {
            "Environment": "INT",
            "StoppedByCapuchin": "true"
        },
        "Tag": "StoppedByCapuchin"
    },
    "Terminate": {
        "Instances": {
            "Environment": "UAT",
            "Capuchin": "eligible-to-destroy"
        },
        "Autoscaling": {
            "Environment": "UAT",
            "Capuchin": "eligible-to-destroy"
        },
        "Tag": "TerminatedByCapuchin"
    },
    "Stop": {
        "Instances": {
            "Environment": "INT",
            "Capuchin": "eligible-to-stop"
        },
        "Tag": "StoppedByCapuchin"
    }
}
```

Using those settings Capuchin will:
- Find all the instances with tags: `Environment = INT` AND `StoppedByCapuchin = true` and start them
removing tag `StoppedByCapuchin`
- Pick one of autoscaling groups with tags: `Environment = UAT` AND `Capuchin = eligible-to-destroy`, 
pick one of its instances there with tags: `Environment = UAT` AND `Capuchin = eligible-to-destroy` and TERMINATE it
tagging it with `TerminatedByCapuchin = true`
- Pick one of instances with tags: `Environment = INT` AND `Capuchin = eligible-to-stop` and STOP IT
tagging it with `StoppedByCapuchin = true`

# License

##### The MIT License (MIT)

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
