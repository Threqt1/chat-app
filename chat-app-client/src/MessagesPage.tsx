import { Divider, InputGroup, InputLeftElement, Icon, Input, InputRightElement, Button, Flex } from "@chakra-ui/react"
import { RiMessage2Line } from "react-icons/ri"
import Message from "./components/Message"
import { useLoaderData } from "react-router-dom"
import { useEffect, useReducer, useRef, useState } from "react"
import { TMessage } from "./types"
import { WebsocketLoaderResult } from "./main"

const MAX_MESSAGES = 50;

export const MessagesPage = () => {
    const { ws, channel, messages } = useLoaderData() as WebsocketLoaderResult & { channel: string, messages: TMessage[] };
    const [chats, updateChats] = useReducer<(chats: TMessage[], update: CHAT_EVENT) => TMessage[], null>(chatReducer, null, () => messages)

    const bottomRef = useRef<HTMLDivElement>(null)

    useEffect(() => {
        ws.setChannel(channel);
        const chatCallback = (chat: TMessage) => {
            updateChats({ type: "add_chat", chat, placeholder: false })
        }
        ws.getEmitter().on("chat_message", chatCallback)

        return () => {
            ws.setChannel("");
            ws.getEmitter().off("chat_message", chatCallback)
        }
    }, [ws, channel])

    useEffect(() => {
        bottomRef.current?.scrollIntoView({ behavior: "instant" })
    }, [chats])

    const onSubmit = (text: string) => {
        if (text.length == 0) return
        updateChats({ type: "add_chat", chat: { id: "temporary", timestamp: -1, from: ws.getUsername(), channelID: channel, content: text }, placeholder: true })
        ws.sendChat(ws.getUsername(), text, channel);
    }

    return <>
        <Flex py={"15px"} ml={"10px"} flex={1} w={"100%"} overflow={"scroll"} direction={"column"} gap={"20px"}>
            {chats.map(chat => <Message name={chat.from} content={chat.content} pending={chat.timestamp < 0} timestamp={chat.timestamp} key={chat.id} />)}
            <div ref={bottomRef} />
        </Flex>
        <Divider borderWidth={"1px"} />
        <MessagesPageInput onSubmit={onSubmit} />
    </>
}

interface MessagesPageInputProps {
    onSubmit: (text: string) => void
}

const MessagesPageInput = (props: MessagesPageInputProps) => {
    const [text, setText] = useState<string>("")

    const onClick = () => {
        const sendText = text.trim()
        setText("")
        props.onSubmit(sendText)
    }

    return <InputGroup w="80%" my="10px">
        <InputLeftElement pointerEvents={"none"}>
            <Icon as={RiMessage2Line} boxSize={7} color={"gray.500"} />
        </InputLeftElement>
        <Input boxShadow={"xl"} placeholder="Enter Message" value={text} onKeyDown={(event) => {
            if (event.key.toLowerCase().trim() === "enter") onClick()
        }} onChange={(event) => setText(event.target.value)} />
        <InputRightElement width="4.5rem">
            <Button size='sm' onClick={onClick}>Send</Button>
        </InputRightElement>
    </InputGroup>
}

interface ADD_CHAT_EVENT {
    type: "add_chat",
    chat: TMessage
    placeholder?: boolean
}

type CHAT_EVENT = ADD_CHAT_EVENT

function chatReducer(
    chats: TMessage[],
    update: CHAT_EVENT
): TMessage[] {
    let newChats = [...chats]
    switch (update.type) {
        case "add_chat":
            newChats.push(update.chat)
            if (!update.placeholder) {
                newChats = newChats.filter(c => !(c.timestamp < 0))
            }
            break;
    }
    return newChats.slice(Math.max(0, newChats.length - MAX_MESSAGES), newChats.length)
}