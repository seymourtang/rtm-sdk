package rtm

import (
	"errors"
	"fmt"
	"time"

	"k8s.io/klog/v2"

	"agora.io/rtm-sdk/internal/rtmlib"
)

type OperatorOptions struct {
	Token  string
	AppID  string
	UserID string
}

type messagePack struct {
	message string
	peer    string
}

type operator struct {
	RtmService             rtmlib.IRtmService
	RtmServiceEventHandler rtmlib.IRtmServiceEventHandler
	receivedCh             chan *messagePack
	sendCh                 chan string
}

const (
	receivedChBuf = 1024
	sendChBuf     = 1024
)

func New(o *OperatorOptions) *operator {
	op := &operator{
		receivedCh: make(chan *messagePack, receivedChBuf),
		sendCh:     make(chan string, sendChBuf),
	}
	rtmService := rtmlib.CreateRtmService()
	klog.Infof("RTM SDK version:%s", rtmlib.GetRtmSdkVersion())
	eventHandler := rtmlib.NewDirectorIRtmServiceEventHandler(op)
	op.RtmServiceEventHandler = eventHandler
	rtmService.Initialize(o.AppID, newRtmServiceEventHandlerImpl(eventHandler))
	rtmService.Login(o.Token, o.UserID)
	return op
}

func (op *operator) Run(stop <-chan struct{}) error {
	if err := op.ReceivedMessage(stop); err != nil {
		return fmt.Errorf("run err:%s", err.Error())
	}
	klog.Info("exiting...")
	return nil
}

func (op *operator) ReceivedMessage(stop <-chan struct{}) error {
	for {
		select {
		case <-stop:
			klog.Info("rtm transformer hub is shutting down")
			op.shutdown()
			return nil
		case data, ok := <-op.receivedCh:
			if !ok {
				return errors.New("err:received channel closed")
			}
			go op.handle(data)
		}
	}
}

func (op *operator) handle(msg *messagePack) {
	t := time.Now().Unix()
	reply := fmt.Sprintf("received:%s,response:%d", msg.message, t)
	op.SendMessageToPeer(msg.peer, reply)
}

func (op *operator) SendMessageToPeer(to, message string) {
	reply := op.RtmService.CreateMessage()
	reply.SetText(message)
	op.RtmService.SendMessageToPeer(to, reply)
	reply.Release()
}

func (op *operator) shutdown() {
	op.RtmService.Logout()
	op.RtmService.Release()
}

func (op *operator) OnLoginSuccess() {
	klog.Info("Login Success")
}

func (op *operator) OnLoginFailure(arg2 rtmlib.AgoraRtmLOGIN_ERR_CODE) {
	klog.Infof("Login LoginFailure,errCode:%d", arg2)

}

func (op *operator) OnRenewTokenResult(arg2 string, arg3 rtmlib.AgoraRtmRENEW_TOKEN_ERR_CODE) {
	klog.Errorf("Login RenewTokenResult,errCode:%d", arg3)
}

func (op *operator) OnTokenExpired() {
	klog.Error("Login TokenExpired")
}

func (op *operator) OnLogout(arg2 rtmlib.AgoraRtmLOGOUT_ERR_CODE) {
	klog.Errorf("OnLogout,errCode:%d", arg2)
}

func (op *operator) OnConnectionStateChanged(arg2 rtmlib.AgoraRtmCONNECTION_STATE, arg3 rtmlib.AgoraRtmCONNECTION_CHANGE_REASON) {
	klog.Infof("OnConnectionStateChanged,state:%s,reason:%s", arg2, arg3)
}

func (op *operator) OnSendMessageResult(arg2 int64, arg3 rtmlib.AgoraRtmPEER_MESSAGE_ERR_CODE) {
	if arg3 != 0 {
		klog.Errorf("failed to send messageID:%d,errCode:%d", arg2, arg3)
	} else {
		klog.Infof("send message id:%d successfully", arg2)
	}
}

func (op *operator) OnMessageReceivedFromPeer(arg2 string, arg3 rtmlib.IMessage) {
	klog.Infof("received message,peerID:%s,data:%s", arg2, arg3.GetText())
	select {
	case op.receivedCh <- &messagePack{
		message: arg3.GetText(),
		peer:    arg2,
	}:
	default:
		klog.Error("receivedChannel is full")
	}
}

type rtmServiceEventHandlerImpl struct {
	rtmlib.IRtmServiceEventHandler
}

func newRtmServiceEventHandlerImpl(eventHandler rtmlib.IRtmServiceEventHandler) *rtmServiceEventHandlerImpl {
	return &rtmServiceEventHandlerImpl{IRtmServiceEventHandler: eventHandler}
}
