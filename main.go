package main

import "fmt"
import "math"
import "os"
import "encoding/json"
import "flag"
import "io"

type ModeloE struct {
	SLR    float64
	RLR    float64
	STMR   float64
	LSTR   float64
	Ds     float64
	TELR   float64
	WEPL   float64
	T      float64
	Tr     float64
	Ta     float64
	Qdu    float64
	Ie     float64
	Bpl    float64
	Ppl    float64
	BurstR float64
	Nc     float64
	Nfor   float64
	Ps     float64
	Pr     float64
	A      float64
}

var filepath string
var useStdin bool
var verbose bool

func init() {
	flag.BoolVar(&verbose, "v", false, "verbose")
	flag.BoolVar(&useStdin, "stdin", false, "read data from stdin (json)")
	flag.StringVar(&filepath, "f", "", "the input file with the initial data (json)")
	flag.Parse()
}

func main() {
	var file io.Reader
	var err error
	if useStdin {
		file = os.Stdin
	} else if filepath != "" {
		file, err = os.Open(filepath)
		if err != nil {
			panic(err)
		}
	} else {
		flag.Usage()
		return
	}

	var data ModeloE
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		panic(err)
	}
	if verbose {
		dataStr, _ := json.MarshalIndent(data, "", "\t")
		fmt.Println("Datos iniciales:", string(dataStr))
	}
	fmt.Printf("R = %f\n", data.Compute())
}

func log(l float64) float64 {
	return math.Log10(l)
}

func R(Ro, Is, Id, Ie_eff, A float64) float64 {
	return Ro - Is - Id - Ie_eff + A
}

func Ro(SLR, No float64) float64 {
	return 15.0 - 1.5*(SLR+No)
}

func No(Nc, Nos, Nor, Nfo float64) float64 {
	t1 := math.Pow(10, Nc/10)
	t2 := math.Pow(10, Nos/10)
	t3 := math.Pow(10, Nor/10)
	t4 := math.Pow(10, Nfo/10)
	return 10 * log(t1+t2+t3+t4)
}

func Nos(Ps, SLR, Ds, OLR float64) float64 {
	term := Ps - OLR - Ds - 14
	return Ps - SLR - Ds - 100 + 0.004*term*term
}

func Nor(RLR, Pre float64) float64 {
	t1 := Pre - 35
	return RLR - 121 + Pre + 0.008*t1*t1
}

func Pre(Pr, LSTR float64) float64 {
	p := math.Pow(10, (10-LSTR)/10)
	return Pr + 10*log(1+p)
}

func Nfo(Nfor, RLR float64) float64 {
	return Nfor + RLR
}

func Is(Iolr, Ist, Iq float64) float64 {
	return Iolr + Ist + Iq
}

func Iolr(Xorl float64) float64 {
	t1 := 1 + math.Pow(Xorl/8, 8)
	return 20 * (math.Pow(t1, 1.0/8.0) - (Xorl / 8))
}

func Xorl(OLR, No, RLR float64) float64 {
	return OLR + 0.2*(64+No-RLR)
}

func Ist(STMRo float64) float64 {
	t1 := aux_ist(STMRo, -13, 6, 8)
	t2 := aux_ist(STMRo, 1, 19.4, 35)
	t3 := aux_ist(STMRo, -3, 33, 13)

	return 12*t1 - 28*t2 - 13*t3 + 29
}

func aux_ist(STMRo, a, b, c float64) float64 {
	t := 1 + math.Pow((STMRo+a)/b, c)
	return math.Pow(t, 1/c)
}

func STMRo(STMR, T, TELR float64) float64 {
	t1 := math.Pow(10, -STMR/10)
	t2 := math.Exp(-T / 4)
	t3 := math.Pow(10, -TELR/10)

	return -10 * log(t1+t2*t3)
}

func Iq(Y, Z float64) float64 {
	t1 := math.Pow(10, Y)
	t2 := math.Pow(10, Z)
	return 15.0 * log(1.0+t1+t2)
}

func Y(Ro, G float64) float64 {
	return ((Ro - 100) / 15) + (46.0 / 8.4) - (G / 9)
}

func Z(G float64) float64 {
	return (46.0 / 30.0) - (G / 40)
}

func G(Q float64) float64 {
	return 1.07 + 0.258*Q + 0.0602*Q*Q
}

func Q(qdu float64) float64 {
	return 37.0 - 15.0*log(qdu)
}

func Id(Idte, Idle, Idd float64) float64 {
	return Idte + Idle + Idd
}

func Idte(Roe, Re, T, STMR, Ist float64) float64 {
	t1 := Roe - Re
	t2 := math.Sqrt(((t1 * t1) / 4) + 100)
	t3 := 1 - math.Exp(-T)
	Idte := (t1/2 + t2 - 1) * t3
	if STMR > 20 {
		return math.Sqrt(Idte*Idte + Ist*Ist)
	}
	return Idte
}

func Roe(No, RLR float64) float64 {
	return -1.5 * (No - RLR)
}

func Re(TERV float64) float64 {
	return 80 + 2.5*(TERV-14)
}

func TERV(TELR, T, STMR, Ist float64) float64 {
	t1 := 1 + (T / 10)
	t2 := 1 + (T / 150)
	t3 := 6 * math.Exp(-0.3*T*T)

	TERV := TELR - 40*log(t1/t2) + t3
	if STMR < 9 {
		return TERV + (Ist / 2)
	}
	return TERV
}

func Idle(Ro, Rle float64) float64 {
	t1 := Ro - Rle
	return t1/2 + math.Sqrt(((t1*t1)/4)+169)
}

func Rle(WEPL, Tr float64) float64 {
	return 10.5 * (WEPL + 7) * math.Pow(Tr+1, -0.25)
}

func OLR(SLR, RLR float64) float64 {
	return SLR + RLR
}

func Idd(Ta, X float64) float64 {
	if Ta <= 100 {
		return 0
	}
	t1 := math.Pow(X, 6)
	t2 := math.Pow(X/3, 6)
	return 25 * (math.Pow(1+t1, 1/6.0) - 3*math.Pow(1+t2, 1.0/6.0) + 2.0)
}

func X(Ta float64) float64 {
	return log(Ta/100) / log(2)
}

func Ie_eff(Ie, Ppl, BurstR, Bpl float64) float64 {
	return Ie + (95-Ie)*(Ppl/((Ppl/BurstR)+Bpl))
}

func (m ModeloE) Compute() float64 {
	mySTMRo := STMRo(m.STMR, m.T, m.TELR)
	myIst := Ist(mySTMRo)

	myOLR := OLR(m.SLR, m.RLR)
	myNos := Nos(m.Ps, m.SLR, m.Ds, myOLR)
	myNfo := Nfo(m.Nfor, m.RLR)
	myPre := Pre(m.Pr, m.LSTR)
	myNor := Nor(m.RLR, myPre)
	myNo := No(m.Nc, myNos, myNor, myNfo)
	myRo := Ro(m.SLR, myNo)
	myQ := Q(m.Qdu)
	myG := G(myQ)
	myY := Y(myRo, myG)
	myZ := Z(myG)
	myIq := Iq(myY, myZ)
	myXorl := Xorl(myOLR, myNo, m.RLR)
	myIolr := Iolr(myXorl)
	myIs := Is(myIolr, myIst, myIq)

	myRle := Rle(m.WEPL, m.Tr)
	myRoe := Roe(myNo, m.RLR)
	myTERV := TERV(m.TELR, m.T, m.STMR, myIst)
	myRe := Re(myTERV)
	myIdte := Idte(myRoe, myRe, m.T, m.STMR, myIst)
	myX := X(m.Ta)
	myIdd := Idd(m.Ta, myX)
	myIdle := Idle(myRo, myRle)
	myId := Id(myIdte, myIdle, myIdd)
	myIe_eff := Ie_eff(m.Ie, m.Ppl, m.BurstR, m.Bpl)

	return R(myRo, myIs, myId, myIe_eff, m.A)
}
