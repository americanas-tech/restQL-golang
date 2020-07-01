# Game Resources

### heroes 
List all heroes or an specific hero

URL: /heroes/:name

| Parameter        | Description|
| ------------- |-------------| 
| name (optional)     | hero name |

### weapons  
List all weapons or an specific weapon

URL: /weapons/:name

| Parameter        | Description|
| ------------- |-------------| 
| name (optional)      | weapon name |

### actions
Get your hero actions

URL: /actions

| Parameter        | Description|
| ------------- |-------------| 
| weaponId     | ID of the weapon |
| heroId     | ID of the hero |

### bosses  
List all bosses or an specific boss

URL: /bosses/:name

| Parameter        | Description|
| ------------- |-------------| 
| name (optional)     | Boss name |

### bossfight
Get a fight result

URL: /boss/fight

| Parameter        | Description|
| ------------- |-------------| 
| weaponId     | ID of the weapon |
| heroId     | ID of the hero |
| bossId     | ID of the boss |
| actionIds     | List of IDs of the actions |