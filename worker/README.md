MercadoPago CoreServices Export Process
=

Skeleton library for exports proccess

Instalation
---

```
$ go get github.com/mercadolibre/coreservices-team/worker
```

Usage
--- 

Para utilizar el esqueleto es necesario definir los siguientes componentes:

idfinder: Es un componente encargado the parsear el mensaje de bigq y obtener el id del item del kvs.
process: Es el componente encargado de realizar la exportación si la misma se encuentra pendiente.
kvsExport: Es el componente encargado de consultar y actualizar el item de export en kvs
lockclient: Es un componente provisto que hace un wrapper del lock de fury y que se encarga de realizar un lock del item a exportar
objectStorage: Es un componente provisto que hace un wrapper del Object Storage de Fury



Changelog
---

0.0.1 - 2017-12-27 

- Initial commit. 