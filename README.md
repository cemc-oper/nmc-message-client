# nmc-message-client

A message client for NWPC operation systems using NMC monitoring platform.

## Installing

Download the latest release and build source code.

Use `Makefile` to build source code on Linux.

All tools will be installed in `bin` directory.

## Getting started

Use different tools to send messages to different message system.

### nmc_monitor_client

Send message to NMC Monitor Platform using `nmc_monitor_client send` command.

```shell script
nmc_monitor_client \
  production \
  --target "host:port" \
  --source grapesGfs \
  --source-ip "source ip" \
  --type gfsProd \
  --product-interval 3 \
  --file-name gmf.gra.2021042200000.grb2 \
  --absolute-data-name /g2/nwp_pd/NWP_GRAPES_GFS_GMF_POST_DATA/2021042200/togrib2/output/grib2_orig/gmf.gra.2021042200000.grb2 \
  --file-size 354884036 \
  --start-time 2021042200 \
  --forecast-time 000 \
  --status 0 \
  --ignore-error \
  --debug
```

The command above will send message below to NMC Monitor Platform's kafka server with default topic monitor.

```json
{
  "topic": "nwpcproduct",
  "source": "grapesGfs",
  "sourceIP": "10.40.143.29",
  "type": "gfsProd",
  "PID": "gfsProd0008000000038081",
  "ID": "2021042212421559780",
  "datetime": "2021-04-22 12:42:15",
  "fileNames": "gmf.gra.2021042200000.grb2",
  "absoluteDataName": "/g2/nwp_pd/NWP_GRAPES_GFS_GMF_POST_DATA/2021042200/togrib2/output/grib2_orig/gmf.gra.2021042200000.grb2",
  "fileSizes": "354884036",
  "result": 0,
  "resultDesc": {
    "startTime": "2021042200",
    "forecastTime": "000"
  }
}
```

## License

Copyright &copy; 2019-2021, Perilla Roc at nwpc-oper.

`nmc-message-client` is licensed under [MIT License](LICENSE)