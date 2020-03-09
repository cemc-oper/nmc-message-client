# nmc-message-client

A message client for NWPC operation systems.

## Installing

Download the latest release and build source code.

Use `Makefile` to build source code on Linux.

All tools will be installed in `bin` directory.

## Getting started

Use different tools to send messages to different message system.

### nmc_monitor_client

Send message to NMC Monitor Platform using `nmc_monitor_client send` command.

```shell script
nmc_monitor_client send \
	--target 10.20.67.183:9092,10.20.67.216:9092,10.20.67.217:9092 \
	--source nwpc_grapes_gfs \
	--type prod_grib \
	--status 0 \
	--file-name gmf.gra.2019062400006.grb2 \
	--absolute-data-name /g2/nwp_pd/NWP_PST_DATA/GMF_GRAPES_GFS_POST/togrib2/output_togrib2/2019062400/gmf.gra.2019062400006.grb2 \
	--start-time 2019062400 \
	--forecast-time 006 \
	--debug
```

The command above will send message below to NMC Monitor Platform's kafka server (10.20.90.35:9092) 
with default topic monitor.

```json
{
	"source":"nwpc_grapes_gfs",
	"type":"prod_grib",
	"status":"0",
	"datetime":1561363080430,
	"fileName":"gmf.gra.2019062400006.grb2",
	"absoluteDataName":"/g2/nwp_pd/NWP_PST_DATA/GMF_GRAPES_GFS_POST/togrib2/output_togrib2/2019062400/gmf.gra.2019062400006.grb2",
	"desc":"{\"start_time\":\"2019062400\",\"forecast_time\":\"006\"}"
}
```

## License

Copyright &copy; 2019-2020, Perilla Roc at nwpc-oper.

`nmc-message-client` is licensed under [MIT License](LICENSE)