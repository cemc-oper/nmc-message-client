# nmc-monitor-client-go

A Go client for NMC Monitor Platform.

## Installing

Download the latest release and build source code.

Use `Makefile` to build source code on Linux.

`nmc_monitor_client` command will be installed in `bin` directory.

## Getting started

### send

Send message to NMC Monitor Platform using `send` command.

```bash
nmc_monitor_client send \
	--target 10.20.90.35:9092 \
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
	"type":"prod_grib2",
	"status":"0",
	"datetime":1561363080430,
	"fileName":"gmf.gra.2019062400006.grb2",
	"absoluteDataName":"/g2/nwp_pd/NWP_PST_DATA/GMF_GRAPES_GFS_POST/togrib2/output_togrib2/2019062400/gmf.gra.2019062400006.grb2",
	"desc":"{\"start_time\":\"2019062400\",\"forecast_time\":\"006\"}"
}
```

## License

Copyright &copy; 2019, Perilla Roc at nwpc-oper.

`nmc-monitor-client-go` is licensed under [MIT License](LICENSE)