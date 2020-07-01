# Querying Your Resources

So, until now you have made `restQL-http` accessible at `localhost:9000` and your mappings look just like this:

```yaml
mappings:
  launches: "https://api.spacexdata.com/v3/launches"
```

There's nothing special about this YAML file, the procedure is simple, any key you add can be queried with a `from` statement just like this:

```
from $RESOURCE
```

Now let's make a simple query to retrieve all launches. Just do a `POST` HTTP request to `localhost:9000/run-query` with the following body:
```
from launches
```
This is making a `GET` request to the endpoint `https://api.spacexdata.com/v3/launches`, you should see the following result:

```json
{
  "launches": {
    "details": {
      "success": true,
      "status": 200,
      "metadata": {}
    },
    "result": [
      {
        "launch_date_unix": 1143239400,
        "mission_name": "FalconSat",
        "launch_success": false,
        "mission_id": [],
        "is_tentative": false,
        "launch_window": 0,
        "launch_site": {
          "site_id": "kwajalein_atoll",
          "site_name": "Kwajalein Atoll",
          "site_name_long": "Kwajalein Atoll Omelek Island"
        },
        "upcoming": false,
        "tbd": false,
        "details": "Engine failure at 33 seconds and loss of vehicle",
        "launch_failure_details": {
          "time": 33,
          "altitude": null,
          "reason": "merlin engine failure"
        },
        "telemetry": {
          "flight_club": null
        },
        "static_fire_date_utc": "2006-03-17T00:00:00.000Z",
        "tentative_max_precision": "hour",
        "static_fire_date_unix": 1142553600,
        "launch_date_utc": "2006-03-24T22:30:00.000Z",
        "timeline": {
          "webcast_liftoff": 54
        },
...
```
*and it goes on all the way*

Woah, that's a lot of data! There should be a way to guarantee the response will have **only** the data that you need. Got the word? **Only**, this is how you use `only` in restQL:

```
from launches   
    only 
        flight_number
        launch_site.site_name
        mission_name
        links.mission_patch_small
        links.mission_patch
        rocket.rocket_id
        rocket.rocket_name
        rocket.rocket_type
```
This will return:
```json
{
  "launches": {
    "details": {
      "success": true,
      "status": 200,
      "metadata": {}
    },
    "result": [
      {
        "launch_site": {
          "site_name": "Kwajalein Atoll"
        },
        "rocket": {
          "rocket_type": "Merlin A",
          "rocket_id": "falcon1",
          "rocket_name": "Falcon 1"
        },
        "mission_name": "FalconSat",
        "links": {
          "mission_patch": "https://images2.imgbox.com/40/e3/GypSkayF_o.png",
          "mission_patch_small": "https://images2.imgbox.com/3c/0e/T8iJcSN3_o.png"
        },
        "flight_number": 1
      },
      {
        "launch_site": {
          "site_name": "Kwajalein Atoll"
        },
        "rocket": {
          "rocket_type": "Merlin A",
          "rocket_id": "falcon1",
          "rocket_name": "Falcon 1"
        },
        "mission_name": "DemoSat",
        "links": {
          "mission_patch": "https://images2.imgbox.com/be/e7/iNqsqVYM_o.png",
          "mission_patch_small": "https://images2.imgbox.com/4f/e3/I0lkuJ2e_o.png"
        },
        "flight_number": 2
      },
      {
        "launch_site": {
          "site_name": "Kwajalein Atoll"
        },
        "rocket": {
          "rocket_type": "Merlin C",
          "rocket_id": "falcon1",
          "rocket_name": "Falcon 1"
        },
        "mission_name": "Trailblazer",
        "links": {
          "mission_patch": "https://images2.imgbox.com/4b/bd/d8UxLh4q_o.png",
          "mission_patch_small": "https://images2.imgbox.com/3d/86/cnu0pan8_o.png"
        },
        "flight_number": 3
      },
      ...
    ]
  }
}
```
**Note:** for complex objects like `rocket`, you can chain the call to a property, such as `rocket.rocket_name`.

Oh yeah, that's a lot more understandable, and achieving this understandability is as simple as this, no schemas, no extra coding, nothing.

Now let's say you want to search for a specific launch, since it's another resource, you should create another mapping entry at `restql.yml`:

```yaml
mappings:
    launches: "https://api.spacexdata.com/v3/launches"
    oneLaunch: "https://api.spacexdata.com/v3/launches/:id"
```

See the `:id` at the end of the endpoint? This is the name of the parameter you should use the `with` clause with:

```
from oneLaunch 
  with id = 27
  only 
      flight_number
      launch_site.site_name
      mission_name
      links.mission_patch_small
      links.mission_patch
      rocket.rocket_id
      rocket.rocket_name
      rocket.rocket_type
```
This query will perform a `GET` to `https://api.spacexdata.com/v3/launches/27` and you'll see the following as a response:

```json
{
  "oneLaunch": {
    "details": {
      "success": true,
      "status": 200,
      "metadata": {}
    },
    "result": {
      "launch_site": {
        "site_name": "CCAFS SLC 40"
      },
      "rocket": {
        "rocket_type": "FT",
        "rocket_id": "falcon9",
        "rocket_name": "Falcon 9"
      },
      "mission_name": "SES-9",
      "links": {
        "mission_patch": "https://images2.imgbox.com/f6/aa/xDtGo0WJ_o.png",
        "mission_patch_small": "https://images2.imgbox.com/fa/ef/4FBvVReu_o.png"
      },
      "flight_number": 27
    }
  }
}
```
