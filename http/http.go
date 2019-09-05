package http

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"golang.org/x/sync/errgroup"
	"http-client/config"
	"io/ioutil"
	"log"
	chnnproto "morse/http-gateway/grpc/proto"
	"net/http"
	"time"
)

type RequestArgs struct {
	UserID  string
	ChnnID  string
	IndexID string
	Token   string
}

func HandleRequest(apiIndex int, args RequestArgs, num int) {
	var data []byte

	switch apiIndex {
	case 1:
		data, _ = sessionList(args)
	case 2:
		data, _ = offlineMsgList(args)
	case 3:
		data, _ = specificMsgs(args)
	case 4:
		data, _ = updateReadIndex(args)
	case 5:
		data, _ = getReadIndex(args)
	case 6:
		data, _ = getOprMsgList(args)
	case 7:
		data, _ = verifyOprMsgs(args)
	}

	responseData(data, args, num, apiIndex)
}

func sessionList(args RequestArgs) ([]byte, error) {
	lsChnnId := chnnproto.PullChnnInfo{
		SChnnId:    args.ChnnID,
		SReadIndex: args.IndexID,
	}
	pb := chnnproto.PullChnnSessionReq{
		SUserId:    args.UserID,
		LsChnnInfo: []*chnnproto.PullChnnInfo{&lsChnnId},
	}

	data, err := proto.Marshal(&pb)
	if err != nil {
		log.Fatalln("PbMarshal Failed sessionList:", err)
		return nil, err
	}

	return data, nil
}

func offlineMsgList(args RequestArgs) ([]byte, error) {
	message := chnnproto.PullChnnMessageInfo{
		SChnnId:    args.ChnnID,
		SBaseIndex: args.IndexID,
		UDirection: *proto.Uint32(1),
		UCount:     *proto.Uint32(10),
	}

	var lspull []*chnnproto.PullChnnMessageInfo
	lspull = append(lspull, &message)

	pb := chnnproto.PullChnnMessageReq{
		SUserId: args.UserID,
		LsPull:  lspull,
	}

	data, err := proto.Marshal(&pb)
	if err != nil {
		log.Fatalln("PbMarshal Failed offlineMsgList:", err)
		return nil, err
	}

	return data, nil
}

func specificMsgs(args RequestArgs) ([]byte, error) {
	specificMessage := chnnproto.PullSpecifiedMessageInfo{
		SChnnId:         args.ChnnID,
		SSpecifiedIndex: args.IndexID,
	}

	var arrMessage []*chnnproto.PullSpecifiedMessageInfo

	arrMessage = append(arrMessage, &specificMessage)

	pb := chnnproto.PullSpecifiedChnnMessageReq{
		SUserId:        args.UserID,
		LsSpecifiedMsg: arrMessage,
	}

	data, err := proto.Marshal(&pb)
	if err != nil {
		log.Fatalln("Mashal data specificMsgs error:", err)
		return nil, err
	}

	return data, nil
}

func updateReadIndex(args RequestArgs) ([]byte, error) {
	readIndexInfo := chnnproto.ReadIndexInfo{
		SChnnId: args.ChnnID,
		SIndex:  args.IndexID,
	}

	pb := chnnproto.UpdateReadIndexReq{
		SUserId: args.UserID,
		LsIndex: []*chnnproto.ReadIndexInfo{&readIndexInfo},
	}

	data, err := proto.Marshal(&pb)
	if err != nil {
		log.Fatalln("Mashal data updateReadIndex error:", err)
		return nil, err
	}

	return data, err
}

func getReadIndex(args RequestArgs) ([]byte, error) {
	pb := chnnproto.GetReadIndexReq{
		SUserId:  args.UserID,
		LsChnnId: []string{args.ChnnID},
	}

	data, err := proto.Marshal(&pb)
	if err != nil {
		log.Fatalln("Mashal data error:", err)
		return nil, err
	}

	return data, nil
}

func getOprMsgList(args RequestArgs) ([]byte, error) {
	pb := chnnproto.PullChnnOprMsgReq{
		SUserId: args.UserID,
	}

	data, err := proto.Marshal(&pb)
	if err != nil {
		log.Fatalln("Mashal data error:", err)
		return nil, err
	}

	return data, nil
}

func verifyOprMsgs(args RequestArgs) ([]byte, error) {
	verifyOprMsg := chnnproto.VerifyOprMsgInfo{
		SChnnId:    args.ChnnID,
		SLastIndex: args.IndexID,
	}

	pb := chnnproto.VerifyChnnOprMsgReq{
		SUserId: args.UserID,
		LsVerify: []*chnnproto.VerifyOprMsgInfo{
			&verifyOprMsg,
		},
	}

	data, err := proto.Marshal(&pb)
	if err != nil {
		log.Fatalln("Mashal data error:", err)
		return nil, err
	}

	return data, nil
}

func responseData(data []byte, args RequestArgs, num int, apiIndex int) {
	client := &http.Client{}
	var g errgroup.Group

	path := chooseURLByIndex(apiIndex)
	
	for i := 1; i <= num; i++ {
		g.Go(func() error {
			req, err := http.NewRequest("POST", "http://"+config.GetConfig().HostPort+path, bytes.NewReader(data))
			req.Header.Add("content-type", "application/x-protobuf;charset=utf-8")
			req.SetBasicAuth(args.UserID, args.Token)

			startTime := time.Now()

			resp, err := client.Do(req)
			if err != nil {
				log.Fatalln(err)
			}

			switch apiIndex {
			case 1:
				spentTime := time.Since(startTime)

				bodyText, err := ioutil.ReadAll(resp.Body)
				s := chnnproto.PullChnnSessionRsp{}

				if err = proto.Unmarshal(bodyText, &s); err != nil {
					log.Fatalln("Unmarshal Pb Failed:", err)
					break
				}

				log.Println("rsp:", s, "PullChnnSessionRsp_num:", i-1, "SpentTime:", spentTime)
			case 2:
				spentTime := time.Since(startTime)

				bodyText, err := ioutil.ReadAll(resp.Body)
				s := chnnproto.PullChnnMessageRsp{}

				if err = proto.Unmarshal(bodyText, &s); err != nil {
					log.Fatalln("Unmarshal Pb Failed:", err)
					break
				}

				log.Println("rsp:", s, "PullChnnMessageRsp_num:", i-1, "SpentTime:", spentTime)
			case 3:
				spentTime := time.Since(startTime)

				bodyText, err := ioutil.ReadAll(resp.Body)
				s := chnnproto.PullSpecifiedChnnMessageRsp{}

				if err = proto.Unmarshal(bodyText, &s); err != nil {
					log.Fatalln("Unmarshal Pb Failed:", err)
					break
				}

				log.Println("rsp:", s, "PullSpecifiedChnnMessageRsp_num:", i-1, "SpentTime:", spentTime)
			case 4:
				spentTime := time.Since(startTime)

				bodyText, err := ioutil.ReadAll(resp.Body)
				s := chnnproto.UpdateReadIndexRsp{}

				if err = proto.Unmarshal(bodyText, &s); err != nil {
					log.Fatalln("Unmarshal Pb Failed:", err)
					break
				}

				log.Println("rsp:", s, "UpdateReadIndexRsp_num:", i-1, "SpentTime:", spentTime)
			case 5:
				spentTime := time.Since(startTime)

				bodyText, err := ioutil.ReadAll(resp.Body)
				s := chnnproto.GetReadIndexRsp{}

				if err = proto.Unmarshal(bodyText, &s); err != nil {
					log.Fatalln("Unmarshal Pb Failed:", err)
					break
				}

				log.Println("rsp:", s, "GetReadIndexRsp_num:", i-1, "SpentTime:", spentTime)
			case 6:
				spentTime := time.Since(startTime)

				bodyText, err := ioutil.ReadAll(resp.Body)
				s := chnnproto.PullChnnOprMsgRsp{}

				if err = proto.Unmarshal(bodyText, &s); err != nil {
					log.Fatalln("Unmarshal Pb Failed:", err)
					break
				}

				log.Println("rsp:", s, "PullChnnOprMsgRsp_num:", i-1, "SpentTime:", spentTime)
			case 7:
				spentTime := time.Since(startTime)

				bodyText, err := ioutil.ReadAll(resp.Body)
				s := chnnproto.VerifyChnnOprMsgRsp{}

				if err = proto.Unmarshal(bodyText, &s); err != nil {
					log.Fatalln("Unmarshal Pb Failed:", err)
					break
				}

				log.Println("rsp:", s, "VerifyChnnOprMsgRsp_num:", i-1, "SpentTime:", spentTime)
			}

			return err
		})
	}

	if err := g.Wait(); err != nil {
		log.Fatalln("ErrGroup failed:", err)
	}
}

func chooseURLByIndex(index int) (result string) {
	switch index {
	case 1:
		result = "/chnn_pull/v1/get/session/list"
	case 2:
		result = "/chnn_pull/v1/get/offlineMsg/list"
	case 3:
		result = "/chnn_pull/v1/get/specificMsgs"
	case 4:
		result = "/chnn_pull/v1/update/readIndex"
	case 5:
		result = "/chnn_pull/v1/get/readIndex"
	case 6:
		result = "/chnn_pull/v1/get/oprMsg/list"
	case 7:
		result = "/chnn_pull/v1/verify/oprMsgs"
	}

	return
}
