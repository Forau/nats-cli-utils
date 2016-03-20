package services_test

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"encoding/json"

	. "github.com/Forau/nats-cli-utils/services"
)

func reverseBytes(a []byte) {
	for left, right := 0, len(a)-1; left < right; left, right = left+1, right-1 {
		a[left], a[right] = a[right], a[left]
	}
}

func readAvailable(r io.Reader) (rb []byte, re error) {
	var buf [256]byte
	for {
		i, err := r.Read(buf[:])
		if i > 0 {
			rb = append(rb, buf[:i]...)
			if i < 256 {
				break // Good enough for now
			}
		} else if err == io.EOF {
			break
		} else {
			re = err
			break
		}
	}
	return
}

type testServiceRequest struct {
	Auth   string
	Action string
	Data   string
}

type testServiceResponse struct {
	Status int
	Resp   string
}

func setupServices(t *testing.T) {
	// Transform to request, and deligate
	RegisterService(func() (*ServiceConfig, interface{}) {
		sc := NewServiceConfig("toReq")
		return sc, func(in string) (string, error) {
			t.Log("toReq input: ", in)
			var req testServiceRequest
			err := json.Unmarshal([]byte(in), &req)
			if err != nil {
				return "", err
			}
			t.Logf("After unmatshal: %+v", req)
			res, err := sc.Next.Invoke(req)
			t.Logf("toReq res: %+v : %+v", res, err)

			b, err := json.Marshal(res)

			return string(b), err
		}
	})

	// Do a service, of some sort
	RegisterService(func() (*ServiceConfig, interface{}) {
		sc := NewServiceConfig("handleReq")
		appStr := sc.Flag("append", "", "Append text to end of response")
		return sc, func(in testServiceRequest) (res testServiceResponse, err error) {
			t.Log("handleReq input: ", in)
			res.Status = len(in.Auth) + len(in.Action) // Just to do something
			res.Resp = in.Data + *appStr
			t.Log("handleReq output: ", res)
			return
		}
	})

	RegisterService(func() (*ServiceConfig, interface{}) {
		return NewServiceConfig("echo"), func(in interface{}) interface{} {
			return in
		}
	})

	// The net service is a _bad_ example, and will be removed or rewritten
	RegisterService(func() (*ServiceConfig, interface{}) {
		fset := NewServiceConfig("net")
		addr := fset.Flag("addr", "", "Destination address for our connection")

		return fset, func(in []byte) ([]byte, error) {
			conn, err := net.Dial("tcp", *addr)
			if err != nil {
				return nil, err
			}
			defer conn.Close()
			//      conn.Write(in)
			fmt.Fprintf(conn, string(in))
			return readAvailable(conn)
		}
	})

}

func TestInvokeEchoService(t *testing.T) {
	setupServices(t)
	service, _ := ServiceRegister.FindService("echo")
	t.Log("Service: ", service)

	testFn := func(in interface{}) {
		res, err := service.Invoke(in)
		t.Logf("Invoked 'echo' with '%+v', and got '%+v'. Err: %+v\n", in, res, err)
		if err != nil || res != in {
			t.Error("Expected input to be same as output")
		}
	}
	testFn("This is a test string")
	testFn(42)
	testFn(struct {
		Data string
		Num  int
	}{"This is a test string", 42})
}

// Tests the use of net-service.  This function needs to be shorter...
func TestInvokeNetService_happy(t *testing.T) {
	setupServices(t)
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal("Oops", err)
	}
	defer l.Close()
	addr := l.Addr()
	t.Log("Address: ", addr)

	service, _ := ServiceRegister.FindService("net", "-addr", addr.String())
	t.Log("Service: ", service)

	testFn := func(in []byte) {
		go func() {
			con, err := l.Accept()
			if err != nil {
				t.Fatal(err)
			}
			t.Log("Got connection: ", con)
			defer con.Close()
			con.SetReadDeadline(time.Now().Add(time.Millisecond * 10))

			ret, err := readAvailable(con)
			t.Log("Read ", len(ret), "b :", string(ret), " with Err: ", err)
			reverseBytes(ret)
			i, _ := con.Write(ret)
			t.Log("Wrote ", i, "b")
		}()
		res, err := service.Invoke(in)
		resb := res.([]byte) // Evil convert
		t.Log("Got ", len(resb), "b: ", string(resb), ": ", err)
		if err != nil {
			t.Error(err)
		}

		reverseBytes(resb)
		if string(resb) != string(in) {
			t.Error("Expected '", len(string(in)), "', but got '", len(string(resb)), "'")
		}
	}

	testFn([]byte("This is just a test"))
	testFn([]byte("This is a really really long test. it needs to be exactly 1024 bytes to make sure that we handle buffers right 降告済南務国上初座役財暮表送合。注大受現情観梨国中諸組分旧負愛。言本近替栄信義生極上検掲効。際転武述郎料作感力務庭図男稿。聞査井半准平全要捕駆千場。板佐動歳家対馬更市残機一。購最見大倍失断全子謙実持車写連成銀傘経各財。択平苦金載天然元化中異話報夜比。療索東多亡時載予天図職東民案話記避。演提止次者志歳企額住更札止団。署民非載惜球変料覚芸部総議影球上鶏。具民患報演甲金訳済社公階送速彼写学。返記規対本委観横戦能爆初問代野。将新紙第綸定目鍋旅乗住期前改院学内断失断全子謙実持車写連成銀傘経各財。択平苦金載天然元化中異話報夜比。療索東多亡時載予天図職東民案話記避。演提止次者志歳企額住更札止団。署民非載惜球and some more"))

	testFn([]byte("This is a really really long test. it needs to be over 1024 bytes to make sure that we handle buffers right" +
		"降告済南務国上初座役財暮表送合。注大受現情観梨国中諸組分旧負愛。言本近替栄信義生極上検掲効。際転武述郎料作感力務庭図男稿。聞査井半准平全要捕駆千場。板佐動歳家対馬更市残機一。購最見大倍果行果彦円使何。創選安朝実応避味作力念媛竹損村米需日楽。路岡形範供体定合無不軽利動。投子改空航強予理豊産以蹊点利戦海検。" +
		"失断全子謙実持車写連成銀傘経各財。択平苦金載天然元化中異話報夜比。療索東多亡時載予天図職東民案話記避。演提止次者志歳企額住更札止団。署民非載惜球変料覚芸部総議影球上鶏。具民患報演甲金訳済社公階送速彼写学。返記規対本委観横戦能爆初問代野。将新紙第綸定目鍋旅乗住期前改院学内断。資法知海戦載見得幸都伊追試週統。Did we get it?"))
}

func TestServiceChain(t *testing.T) {
	setupServices(t)

	quickFindService := func(args ...string) Service {
		s, _ := ServiceRegister.FindService(args...)
		return s
	}

	srv1, _ := ServiceRegister.FindService("toReq", "handleReq", "-append", "TESTING")
	srv2, _ := ServiceRegister.FindService("toReq", "handleReq", "-append", " SRV2")
	testExpect := []struct {
		Service Service
		In      string
		Out     string
		Err     bool
	}{
		// With service1
		{srv1, `{"auth": "my-token", "action": "test", "data": "0123456789"}`,
			`{"Status":12,"Resp":"0123456789TESTING"}`, false},
		{srv1, `{"auth": "my-tok", "action": "test", "data": ""}`,
			`{"Status":10,"Resp":"TESTING"}`, false},
		// And service2
		{srv2, `{"auth": "my-token", "action": "EatCookie", "data": "Santa's little helper"}`,
			`{"Status":17,"Resp":"Santa's little helper SRV2"}`, false},
		// And echo (We actually get change from input, since json to string will capitalize)
		{quickFindService("toReq", "echo"),
			`{"auth": "my-token", "action": "EatCookie", "data": "Santa's little helper"}`,
			`{"Auth":"my-token","Action":"EatCookie","Data":"Santa's little helper"}`, false},
		// And create a service on the fly, to check we wont dissrupt srv1 or srv2
		{quickFindService("toReq", "handleReq", "-append", ". (FlyBoy)"),
			`{"action": "Smoke me a kipper!", "data": "I'll be back for breakfast"}`,
			`{"Status":18,"Resp":"I'll be back for breakfast. (FlyBoy)"}`, false},
		// Retest with srv1
		{srv1, `{"auth": "my-token", "action": "test", "data": "0123456789"}`,
			`{"Status":12,"Resp":"0123456789TESTING"}`, false},
		// Parse error
		{srv2, `{"auth": "my-token", "action": `, "", true},
	}
	for _, te := range testExpect {
		t.Log("Testing: ", te)
		res, err := te.Service.Invoke(te.In)
		t.Log("Result: ", res, err)
		hasErr := err != nil
		if res != te.Out || hasErr != te.Err {
			t.Error("Not as expected (", te.Out, te.Err, ")")
		}
	}

	// Test service not found
	for _, args := range [][]string{
		{"沒有"},
		{"toReq", "沒有"},
		{"echo", "blahablaha..."},
	} {
		t.Logf("Trying to find service '%+v'", args)
		s, e := ServiceRegister.FindService(args...)
		t.Logf("Got response %+v, %+v", s, e)
		if e == nil {
			t.Error("Expected an error")
		}
	}

}
