import React, { useEffect, useRef } from 'react'
// import { useParams } from 'react-router-dom';

const Room = () => {
    const { room_id} = useParams();

    const userVideo = useRef()
    const partnerVideo =useRef()
    const userStream = useRef()
    const peerRef = useRef()
    const WebSocketRef = useRef()

    
        const openCamera = async () =>{
            const allDevices = await navigator.mediaDevices.enumerateDevices()
            const cameras = allDevices.filter((device) => device.kind == "videoinput")
            const constraints = { 
                audio:true,
                video: {
                    deviceId :cameras[1].deviceId,  
              },
            };
        

        try {
            return await navigator.mediaDevices.getUserMedia(constraints)
        } catch (error) {
            console.log(error)
        }
    };

    useEffect(()=>{
        // const ws = new WebSocket(`ws://localhost:8000/join?roomID=${room_id}`)
        // ws.addEventListener("open",()=>{
        //     ws.send(JSON.stringify({join:"true"}));
        // }); 
        openCamera().then((stream)=>{
                userVideo.current.srcObject = stream
                userStream.current = stream
                
                WebSocketRef.new WebSocket{
                    `ws://localhost:8000/join?roomID=${room_id}`
                };
                WebSocketRef.current.addEventListener("open", ()=>{
                    WebSocketRef.current.send(JSON.stringify({join: true}));
                })

                WebSocketRef.current.addEventListener("message",async (e)=>{
                    const message = JSON.parse(e.data)

                    if(message.join) {
                            callUser()
                    }

                    if(message.iceCandidate){
                        console.log("recived and adding ice candidate")
                    
                        try {
                            await peerRef.current.addIceCandidate(message.iceCandidate)
                        } catch (error) {
                            console.log("error reciving ice candidate",error)
                        }
                    }

                    if (message.offer) {
                        handleOffer(message.offer)
                    }

                    if(message.answer){
                        console.log("Recived answer")
                        peerRef.current.setRemoteDescription(
                            new RTCSessionDescription(message.answer)
                        )
                    }
                });

        
            });

    });

    const handleOffer = (offer) => {
        console.log("recived offer creating answer")
        peerRef.current = createPeer();

        await peerRef.current.setRemoteDescription(
            new RTCSessionDescription(offer)
        );
        
        userStream.current.getTracks().forEach((track)=>{
            peerRef.current.addTrack(track, userStream.current)
        })
   
        const answer = await peerRef.current.createAnswer();
        await peerRef.current.setLocalDescription(answer);


        WebSocketRef.current.send(
            JSON.stringify({answer: peerRef.current.localDescription})
        )
    }

const callUser =()=>{
        console.log("calling other user");
        peerRef.current = createPeer();

        userStream.current.getTracks().forEach(track => {
            peerRef.current.addTrack(track,userStream.current)            
        });
    }

const createPeer =()=>{
    console.log("creating peer connection")
    const peer = new RTCPeerConnection({
        iceServers : [{urls: "stun:stun.l.google.com:19302"}],
    });

    peer.onnegotiationneeded= handleNegotiationNeeded;
    peer.onicecandidate = handleIceCandidateEvent;
    peer.ontrack = handleTrackEvent;

    return peer
};

const handleNegotiationNeeded =() =>{
    console.log("creating the offer");
    try {
        const myOffer = await peerRef.current.createOffer();
        await peerRef.current.setLocalDescription(myOffer);

        WebSocketRef.current.send(
            JSON.stringify({offer: peerRef.current.localDescription})
        )

    } catch (error) {
        console.log(error)
    }
}
const handleIceCandidateEvent =(e)=>{
    console.log("found Ice candidate");
    WebSocketRef.current.send(JSON.stringify({iceCandidate: e.candidate}))
}
const handleTrackEvent = (e) =>{
    console.log("recived tracks")
    partnerVideo.current.srcObject = e.stream[0]
}




  return (
    <div>
            <video autoPlay controls={true} ref={userVideo}></video>
            <video autoPlay controls={true} ref={partnerVideo}></video>
    </div>
  )
}

export default Room
