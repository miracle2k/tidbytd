# tidbytd

The [Tidbyt](https://tidbyt.com/) is an incredibly fun device. Creating custom apps
for it is fun and easy. But getting to enjoy your own creations isn't. Getting an
app into the app store can take a while. And of course, you may have private apps that
simple aren't suitable for public release.

This is a very basic, very much WIP rendering daemon which you can run on your own server.
You can point it to one or more .star files, and one or more devices, and it will, in regular
intervals, render the apps and push them to your device via the Tidbyt APIs.

It does *not* serve your Tidbyt device directly.

### Usage

Create a config.yaml file in this format:

```
devices:
  foo:
    id: severely-foobar-ae3
    api_token: "ey..."
content:
- url: /Users/michael/Dev/tidbyt-community/apps/aiclock/ai_clock.star
  name: aiclock
- url: https://github.com/tidbyt/community/blob/d1495b92b71fad9b3190f5592d6a82aceec7dbcf/apps/nouns/nouns.star
  name: nouns
```
   

Run:

   $ ./tidbytd


### Notes

On M1 Macs:

   export LIBRARY_PATH=$LIBRARY_PATH:/usr/local/lib:/opt/homebrew/lib/
   export CPATH=$CPATH:/usr/local/include:/opt/homebrew/include/

See: https://github.com/le0pard/webp-ffi/issues/14#issuecomment-975479834
