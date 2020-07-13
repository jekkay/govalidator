# go-validator

---

<p>A Fast Validator for struct filed using tag, highly inspired by automatically recovery theory.</p>

<p>There are numerous validators for golang, but almost of them lack of automatically 
recovery from error. It's very easy to find errors exist during run-time. But if we wanna our process 
 is to be sure running in safe mode to keep our service is perfect available, it definitely require
 self-recovery ability, especially when wrong setting or arguments can be delivered
 by user, or other programs...
 </p>

## quick start

```

import (
    "github.com/jekkay/govalidator"
    "encoding/json"
    "fmt"
    "testing"
}

type Range struct {
    A int32   `json:"a" min:"10" max:"100" default:"50"`
    B int32   `json:"b" min:"20" max:"90" default:"80"`
    C *uint64 `json:"c" min:"30" max:"90" req:"true" default:"60"`
    D string  `json:"d" min:"1" max:"10" req:"true" in:"hello,world,jekkay" regex:"^[a-d]+$" default:"jekkay"`
}

func TestValidObject2(t *testing.T) {
	r := new(Range)
	r.A = 120
	r.B = 130

	if e := govalidator.ValidObject(r, false); e != nil {
		fmt.Println(e)
	}
	govalidator.ValidObject(r, true)
	bs, _ := json.MarshalIndent(r, "", "  ")
	fmt.Println(string(bs))
}
```

output is 
```
[
`A` at most 100, current is 120
`B` at most 90, current is 130 
`C` is empty
`D` is empty
]
{
  "a": 100,
  "b": 90,
  "c": 60,
  "d": "jekkay"
}
```

<p><code>ValidObject</code>function will check struct filed whether is valid, 
it will return the first error and skip next checking. Set parameter `fix` to <code>true</code> 
if you wanna auto fix struct filed. It's very helpful to make sure the value in custom range, and set
default value if necessary.</p>

it's available for the nested struct as well, see test file for more: [validator_test.go](./validator_test.go).


## Function

| function | parameters | description |
|--------|--------|--------|
| ValidObject | <code>obj</code>:ptr, a pointer to struct object<br/><code>fix</code>: boolean, indicate whether auto adjust value<br/>| auto fix error value|
| Validate | <code>obj</code>:ptr, a pointer to struct object<br/>| validate object error, return the first one |
| Validates | <code>obj</code>:ptr, a pointer to struct object<br/>|  validate object error, return all errors |

## Auto Fix Value


<p>Set parameter `fix` to <code>true</code>, <code>ValidObject</code>function will do some adjustment.</p>

 - if the filed is number(<code>int</code>, <code>uint</code>...), logic like this:

```
   if currentValue < min {
       if currentValue == 0 {
            currentValue = defalut
       } else {
            currentValue = min
       }
   }
   if currentValue > max {
       currentValue = max
   }
     
``` 
<p>if there is no <code>default</code> tag set, <code>min</code> will be used as default value instead. </p>
 
 - if the field is string(<code>string</code>), logic like this

```
  if len(str) < min || len(str) > max 
       || dismatch(str, regex) || !in(str, options) {
      str = default
  }
```

 - if the field is pointer, logic like this:
 
```
  if ptr == nil && required{
    ptr = new(type)
  }
  if ptr != nil {
     check(*ptr) --> logic go above
  }

```

## Tags

<p>Tags are used to describe the constraint of the field.</p>


| Tag | description | require |
|------|------|------|
| default | default value | √ | 
| min | the minimum length of the string,<br/> or the minimum value of number | × |
| max | the maximum length of the string ,<br/> or the maximum value of number| × |
| req | require value | × |
| regex | regular expression | × |
| in | value options | × |

<p>Number Range Constraint:</p>

| Tag | Field Type |description |
|------|------|------|
| min | Int, Int8,Int16,Int32,Int64<br/>Uint,Uint8,Uint16,Uint32,Uint64| minimum value |
| max | Int, Int8,Int16,Int32,Int64<br/>Uint,Uint8,Uint16,Uint32,Uint64| Maximum value |
| default | Int, Int8,Int16,Int32,Int64<br/>Uint,Uint8,Uint16,Uint32,Uint64| defalut value |


| Field Type | min | max | default |
|-------|-------|-------|-------|
| Uint | 0 | 2^32 -1 | 0 |
| Uint8 | 0 | 2^8 -1 | 0 |
| Uint16 | 0 | 2^16 -1 | 0 |
| Uint32 | 0 | 2^32 -1 | 0 |
| Uint64 | 0 | 2^64 -1 | 0 |


| Field Type | min | max | default |
|-------|-------|-------|-------|
| Int | -2^31 | 2^31 -1 | 0 |
| Int8 | -2^7 | 2^7 -1 | 0 |
| Int16 | -2^15 | 2^15 -1 | 0 |
| Int32 | -2^31 | 2^31 -1 | 0 |
| Int64 | -2^63 | 2^63 -1 | 0 |

| Field Type | min | max | default | precision|
|-------|-------|-------|-------|-------|
| Float32 | -3.4e38 | 3.4e38 | 0 | 7 |
| Float64 | -1.7e308 | 1.7e308	 | 0 | 16 |


<p> Tags decorate description:</p>

| Field Type | min | max | default | req | in | regex |
|-------|-------|-------|-------|-------|-------|-------|
| (number)<br/>Int8,Uint8,<br/>Int16,Uint16<br/>...<br/>Int64,Uint64| √ | √ | √ | × |× |× |
| string | √ | √ | √ |  √ | √ | √ |
| ptr -> number |√ | √ | √ | √ |× |× |
| ptr -> string | √ | √ | √ |  √ | √ | √ |
| ptr -> struct |× |× |× |√ |× |× |

## Author
 - Jekkay Hu
 - Blog: [http://www.easysb.cn](http://www.easysb.cn)
 - Email: jekkay#qqvips.cn
 