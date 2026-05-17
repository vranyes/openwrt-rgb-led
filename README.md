# rgb-led (OpenWRT)

This repository contains a small Go program (`main.go`) that performs a smooth RGB LED sweep by writing to the LED sysfs entries. The repository includes a sample OpenWRT `procd`-compatible init script at `files/etc/init.d/rgb-led` (note: the filename in the repo is `rgb-led`); when installed on a router it should be named `/etc/init.d/rgb-led` and will run the compiled binary.

Quick install

1. Build the binary for your router's architecture (example for MIPS little-endian):

    ```sh
    GOOS=linux GOARCH=mipsle go build -o rgb-led main.go
    ```

2. Copy the binary and the init script to your router and enable the service:

```sh
scp rgb-led root@ROUTER:/usr/bin/
scp files/etc/init.d/rgg-led root@ROUTER:/etc/init.d/rgb-led
ssh root@ROUTER 'chmod +x /usr/bin/rgb-led /etc/init.d/rgb-led'
ssh root@ROUTER '/etc/init.d/rgb-led enable'
ssh root@ROUTER '/etc/init.d/rgb-led start'
```

Runtime configuration

- The program reads environment variables for LED names and timing: `RED_LED`, `GREEN_LED`, `BLUE_LED`, `STEP`, `DELAY`.
- You can override defaults by exporting env vars before calling start, or edit `/etc/init.d/rgb-led` defaults.

Examples

```sh
# start with faster steps (STEP=5) and shorter delay (DELAY=200ms)
ssh root@ROUTER 'RED_LED=red:indicator STEP=5 DELAY=200 /etc/init.d/rgb-led start'

# check status
ssh root@ROUTER '/etc/init.d/rgb-led status'

# stop service
ssh root@ROUTER '/etc/init.d/rgb-led stop'
```

UCI configuration

- A sample UCI config file is included at `files/etc/config/rgb-led`; install it to `/etc/config/rgb-led` on the router. The file looks like:

```sh
config rgb_led 'main'
	option enabled '1'
	option red_led 'red:indicator'
	option green_led 'green:indicator'
	option blue_led 'blue:power'
	option step '1'
	option delay '1000'
```

- To write the file on the router (simple):

```sh
cat >/etc/config/rgb-led <<'EOF'
config rgb_led 'main'
	option enabled '1'
	option red_led 'red:indicator'
	option green_led 'green:indicator'
	option blue_led 'blue:power'
	option step '1'
	option delay '1000'
EOF
uci commit rgb-led
```

- Or change values with `uci` (when `main` section exists):

```sh
uci set rgb-led.main.red_led='red:indicator'
uci set rgb-led.main.step='2'
uci commit rgb-led
```

- After changing UCI values, restart the service:

    ```sh
    /etc/init.d/rgb-led restart
    ```

Notes

- Build the binary for the correct CPU architecture of your router; common values: `arm`, `arm64`, `mips`, `mipsle`, `amd64`.
- For packaging into an OpenWRT ipk, include `files/etc/init.d/rgb-led` (the repo filename) in the package payload and ensure the installed path is `/etc/init.d/rgb-led` on the device.
