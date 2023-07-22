# tidbytd

Unofficial tidbyt rendering server. Does not replace the official
Tidbyt server, but rather does background pushes to it.

   $ vim config.yaml
   $ ./tidbytd

### Notes

On M1 Macs:

   export LIBRARY_PATH=$LIBRARY_PATH:/usr/local/lib:/opt/homebrew/lib/
   export CPATH=$CPATH:/usr/local/include:/opt/homebrew/include/

See: https://github.com/le0pard/webp-ffi/issues/14#issuecomment-975479834
