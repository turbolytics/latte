# signals-collector
Collect metrics from a variety of operational sources



## API

### Data Types

**Collector**

```
{
    "name": <str>,
    "metric": {
        "name": "" 
        "type": "GAUGE"
    }
    
    "schedule": {
        "interval": <go time interval> 
        "cron": Cron expression
    },
    
    "collector": {
        type: <MONGODB,postgres,stripe,salesforce,etc>,
        "config": {
            
        }
    },
    
    "output": {
    }
}
```


### Output Type

**Metric**
```
{
    name: <str>
    type: <enum:>
    value: {
        v:
    },
    dimensions: {
        map<str, str> 
    }
}
```



### REST

- `GET /collectors`
 
Lists all collectors.

- `POST /collectors`

Create a new collector.

