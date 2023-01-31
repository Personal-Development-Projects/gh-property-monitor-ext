# gh-property-monitor-ext
A gh cli extension for monitoring property files for changes


## Installing the plugin locally

run the following command once authorized within the GH CLI 

1) run ```git auth status and ensure you are logged in at the appropiate GitHub host```
2) run ```gh ext install Personal-Development-Projects/gh-property-monitor```

## Creating a new release

There is currently a workflow available on this repo that will compile and output the release and its associated assets.

To trigger this action utilize a tag approach as laid out below
  
  1) ```git tag v1.X.X```
  2) ```git push origin v1.X.X ``` Push this new tag up to GitHub
  3) ```gh run view``` Ensure you see the gh workflow run 
  4) Further you can ensure the release was created by running ```gh release view```

## Building new extension to test locally


## Example Use Case Results

Once authorized on your desired remote github server install the extension in one of the options detailed

![GH CLI Auth Status](Documentation/images/GHCliAuth.png?raw=true "GH CLI Auth Status")

This can be installed through cloning the repo, accessing the releases, or utilizing the gh extension search function.

Once installed run it using the ```gh property-monitor```

![Terminal Output](Documentation/images/TerminalOutput.png?raw=true "Terminal Output")

![Excel Outputted File](Documentation/images/ExcelOutput.png?raw=true "Excel Outputted File")
