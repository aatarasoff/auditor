# Auditor

[![Join the chat at https://gitter.im/aatarasoff/auditor](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/aatarasoff/auditor?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Audit logging bridge for Docker

Auditor automatically logs start/stop events of services for Docker containers based on published ports and metadata from the container environment. Based on [Registrator](https://github.com/gliderlabs/registrator) and supports Elasticsearch and Logstash engine.

By default, it can register services without any user-defined metadata. This means it works with *any* container, but allows the container author or Docker operator to override/customize the service definitions.

## Getting Auditor

You can get the latest release of Auditor via Docker Hub:

	$ docker pull aatarasoff/auditor:latest

## Starting Auditor

	$ docker run -d \
		-v /var/run/docker.sock:/tmp/docker.sock \
		-h $HOSTNAME aatarasoff/auditor <engine-uri>

### Logstash Engine (recommended)

To use the Logstash, specify a Logstash TCP host and port. If no host is provided, `127.0.0.1:5959` is used. Examples:

	$ auditor logstash://10.0.0.1:5959
	$ auditor logstash:

### Elasticsearch Engine

To use the Elasticsearch directly, specify a ES host and port. If no host is provided, `127.0.0.1:9200` is used. Examples:

	$ auditor elastic://10.0.0.1:9200
	$ auditor elastic:

Also you need to add logstash-like template to ES:
```
{
  "template" : "containers*",
  "settings" : {
    "index.refresh_interval" : "5s"
  },
  "mappings" : {
    "_default_" : {
       "_all" : {"enabled" : true},
       "dynamic_templates" : [ {
         "string_fields" : {
           "match" : "*",
           "match_mapping_type" : "string",
           "mapping" : {
             "type" : "string", "index" : "analyzed", "omit_norms" : true,
               "fields" : {
                 "raw" : {"type": "string", "index" : "not_analyzed", "ignore_above" : 256}
               }
           }
         }
       } ],
        "_ttl": {
         "enabled": true,
         "default": "1d"
       },
       "properties" : {
         "@version": { "type": "string", "index": "not_analyzed" },
         "geoip"  : {
           "type" : "object",
             "dynamic": true,
             "path": "full",
             "properties" : {
               "location" : { "type" : "geo_point" }
             }
         }
       }
    }
  }
}
```
