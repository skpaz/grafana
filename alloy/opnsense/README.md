# OPNsense & Grafana Alloy

>[!NOTE]
> This is a work in progress.

## Install Alloy

```bash
tmpfile=$(mktemp)
installdir="/usr/local/sbin"
curl -L https://github.com/grafana/alloy/releases/download/FIXME_VERSION_NUMBER/alloy-freebsd-amd64.zip -o ${tmpfile}
unzip -d ${installdir} ${tmpfile}
mv ${installdir}/alloy-freebsd-amd64 ${installdir}/alloy
chmod +x ${installdir}/alloy
rm -rf ${tmpfile}
unset installdir tmpfile
```

Some Alloy components require root, so `/usr/local/sbin` is the most appropriate location.

## Configure Alloy

```bash
sudo sh -c 'cat > /usr/local/etc/alloy/config.alloy' << 'EOF'

FIXME_ADD_CONFIG_ALLOY

EOF
```

## Set up Service

Create service:

```bash
sudo sh -c 'cat > /usr/local/etc/rc.d/alloy' << 'EOF'
#!/bin/sh

# PROVIDE: alloy
# REQUIRE: NETWORKING
# KEYWORD: shutdown

#
# Add this to /etc/rc.conf.local to enable this service:
#
# alloy_enable="YES"
#

. /etc/rc.subr

name="alloy"
desc="Grafana Alloy telemetry agent"
rcvar=alloy_enable

: ${alloy_enable="NO"}

pidfile="/var/run/${name}.pid"
config_file="/usr/local/etc/${name}/config.alloy"
output_file="/var/log/${name}.log"
required_files="${config_file}"

exec_path="/usr/local/sbin/${name}"
exec_args="run ${config_file} --server.http.listen-addr FIXME_HTTP_LISTEN_IP:12345"

command="/usr/sbin/daemon"
command_args="-t ${name} -o ${output_file} -P ${pidfile} ${exec_path} ${exec_args}"

load_rc_config ${name}
run_rc_command "$1"

EOF
```

Set execution bit:

```bash
chmod +x /usr/local/etc/rc.d/alloy
```

Enable the service on boot:

```bash
sudo sh -c 'cat >> /etc/rc.conf.local' << 'EOF'
alloy_enable="YES"
EOF
```

## Test Service

Start the service to test it:

```bash
service alloy start
```

This will start the service in interactive mode. You will see the `stdout` and `stderr` output in the terminal, which will allow you to easily troubleshoot any issues with your Alloy configuration file, if needed.

Once Alloy is confirmed functional, hit Ctrl-C to stop the service.

## Reboot

In order to start Alloy in non-interactive mode, reboot the system.
