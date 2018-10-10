install the service file:

```
ln -s systemd.service /etc/systemd/system/caserver.service
```

enable the service

```
systemctl enable nginx_controller.service
```

start the service

```
systemctl start nginx_controller.service
```