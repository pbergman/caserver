install the service file:

```
ln -s systemd.service /etc/systemd/system/caserver.service
```

enable the service

```
systemctl enable caserver.service
```

start the service

```
systemctl start caserver.service
```
