# go-validator
A simple validator for struct filed using tag

---

## quick start

```
type Range struct {
	A int32   `json:"a" min:"10" max:"100" default:"50"`
	B int32   `json:"b" min:"20" max:"90" default:"80"`
	C *uint64 `json:"c" min:"30" max:"90" default:"60"`
}

func TestValidObject2(t *testing.T) {
	r := new(Range)
	r.A = 120
	r.B = 130
	r.C = new(uint64)
	*r.C = 0

	if e := ValidObject(r, false); e != nil {
		fmt.Println(e)
	}
	ValidObject(r, true)
	bs, _ := json.MarshalIndent(r, "", "  ")
	fmt.Println(string(bs))
}
```

output is 
```
A Maximum value is 100，current is 120
{
  "a": 100,
  "b": 90,
  "c": 60
}
```

<p><code>ValidObject</code>function will check struct filed whether is valid, 
it will return the first error and skip next checking. Set parameter `fix` to <code>true</code> 
if you wanna auto fix struct filed. It's very helpful to make sure the value in custom range, and set
default value if necessary.</p>


## Auto Fix Value

<p>Set parameter `fix` to <code>true</code>, <code>ValidObject</code>function will do some adjustment.</p>

 - if the filed is number(<code>int</code>, <code>uint</code>...), logic is like this:

```
   if currentValue < min {
       if currentValue == 0 {
            currentValue = defalutValue
       } else {
            currentValue = min
       }
   }
   if currentValue > max {
       currentValue = max
   }
     
``` 

<p>if there is no <code>defalut</code> tag set, <code>min</code> will be used as default value instead. </p>

## Tag

<p>Tags are used for describe the constraint of the field.</p>

| Tag | Field Type |description |
|------|------|------|
| min | Int, Int8,Int16,Int32,Int64<br/>Uint,Uint8,Uint16,Uint32,Uint64| minimum value |
| max | Int, Int8,Int16,Int32,Int64<br/>Uint,Uint8,Uint16,Uint32,Uint64| Maximum value |
| default | Int, Int8,Int16,Int32,Int64<br/>Uint,Uint8,Uint16,Uint32,Uint64| defalut value |


## TO BE CONTINUED

## Author
 - Jekkay Hu, Blog: [http://www.eassyb.cn](http://www.eassyb.cn)
 