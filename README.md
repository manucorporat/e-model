# G.107 E-model implementation in Go

This is an Open source implementation in Go of [G.107 : The E-model: a computational model for use in transmission planning](https://www.itu.int/rec/dologin_pub.asp?lang=e&id=T-REC-G.107-201506-I!!PDF-E&type=items).



## Usage

Given a json file with all the initial values of the model like this:
```
cat modelo.json
```
```json
{
  "SLR"    : 8,
  "RLR"    : 2,
  "STMR"   : 15,
  "LSTR"   : 18,
  "Ds"     : 3,
  "TELR"   : 65,
  "WEPL"   : 110,
  "T"      : 0,
  "Tr"     : 0,
  "Ta"     : 0,
  "Qdu"    : 1,
  "Ie"     : 0,
  "Bpl"    : 1,
  "Ppl"    : 0,
  "BurstR" : 1,
  "Nc"     : -70,
  "Nfor"   : -64,
  "Ps"     : 35,
  "Pr"     : 35,
  "A"      : 0
}
```


```
➜  modelo_e go run main.go -f modelo.json
R = 93.206208
```