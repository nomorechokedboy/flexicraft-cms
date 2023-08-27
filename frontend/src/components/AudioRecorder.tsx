import { createSignal, onCleanup } from "solid-js";

const MIME_TYPE = "audio/webm";

export default function AudioRecorder() {
    const [permission, setPermission] = createSignal(false);
    const [stream, setStream] = createSignal<MediaStream>();
    let mediaRecorder: MediaRecorder;
    const [recordingStatus, setRecordingStatus] = createSignal("inactive");
    const [audioChunks, setAudioChunks] = createSignal<Blob[]>([]);
    const [audio, setAudio] = createSignal<string[]>([]);

    const startRecording = async () => {
        setRecordingStatus("recording");
        //create new Media recorder instance using the stream
        const s = stream();
        if (!s) {
            return;
        }

        const media = new MediaRecorder(s, { mimeType: MIME_TYPE });
        //set the MediaRecorder instance to the mediaRecorder ref
        mediaRecorder = media;
        //invokes the start method to start the recording process
        mediaRecorder.start();
        let localAudioChunks: Blob[] = [];
        mediaRecorder.ondataavailable = (event) => {
            if (typeof event.data === "undefined") return;
            if (event.data.size === 0) return;
            localAudioChunks.push(event.data);
        };
        setAudioChunks(localAudioChunks);
    };

    function cleanupObjectUrl() {
        audio().forEach((url) => {
            console.log("Cleanup url: ", { url });

            URL.revokeObjectURL(url);
        });
    }

    const stopRecording = () => {
        // cleanupObjectUrl();
        setRecordingStatus("inactive");
        //stops the recording instance
        mediaRecorder.stop();
        mediaRecorder.onstop = () => {
            //creates a blob file from the audiochunks data
            const audioBlob = new Blob(audioChunks(), { type: MIME_TYPE });
            //creates a playable URL from the blob file.
            const audioUrl = URL.createObjectURL(audioBlob);
            audio().push(audioUrl);
            setAudio(audio());
            console.log({ audioUrl, audio: audio() });

            setAudioChunks([]);
        };
    };

    const getMicrophonePermission = async () => {
        if ("MediaRecorder" in window) {
            try {
                const streamData = await navigator.mediaDevices.getUserMedia({
                    audio: true,
                    video: false,
                });
                setPermission(true);
                setStream(streamData);
            } catch (e) {
                console.error(e);
            }
        } else {
            alert("The MediaRecorder API is not supported in your browser.");
        }
    };

    onCleanup(() => {
        console.log("Cleanup running...");
        cleanupObjectUrl();
    });

    return (
        <div>
            <h2>Audio Recorder</h2>
            <main>
                <div class="audio-controls">
                    {!permission() ? (
                        <button onClick={getMicrophonePermission} type="button">
                            Get Microphone
                        </button>
                    ) : null}
                    {permission() && recordingStatus() === "inactive" ? (
                        <button onClick={startRecording} type="button">
                            Start Recording
                        </button>
                    ) : null}
                    {recordingStatus() === "recording" ? (
                        <button onClick={stopRecording} type="button">
                            Stop Recording
                        </button>
                    ) : null}{" "}
                </div>
                {audio ? (
                    <div class="audio-container">
                        <audio src={audio()} controls></audio>
                        <a download href={audio()}>
                            Download Recording
                        </a>
                    </div>
                ) : null}
            </main>
        </div>
    );
}
