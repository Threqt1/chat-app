import { Box, Button, Divider, HStack, Input, Wrap, useToast } from "@chakra-ui/react"
import { useLoaderData } from "react-router-dom"
import Channel from "./components/Channel"
import { TChannel } from "./types"
import { useEffect, useState } from "react"
import { CreateChannel, ErrorUsernameDoesntExist, GetChannel } from "./APIWrapper"
import { WebsocketLoaderResult } from "./main"

export const ChannelsPage = () => {
    const { ws, ...loaderData } = useLoaderData() as WebsocketLoaderResult & { channels: { channels: TChannel[] } }
    const errorToast = useToast()
    const [channels, setChannels] = useState<TChannel[]>(loaderData.channels.channels)

    useEffect(() => {
        const channelAddedCallback = async (channelID: string) => {
            const channel = await GetChannel(channelID)
            if (!channel.data) return
            setChannels((curr) => [...curr, channel.data!.channel])
        }
        ws.getEmitter().on("channel_added", channelAddedCallback)

        return () => {
            ws.getEmitter().off("channel_added", channelAddedCallback)
        }
    }, [ws])

    const invalidInput = (type: string, error: string) => {
        errorToast({
            title: `Failed to ${type}.`,
            description: error,
            status: "error",
            duration: 2000,
            isClosable: true
        })
    }

    const onSubmit = async (text: string) => {
        if (text.length === 0) return
        let errorMsg = "An unknown error occurred."
        try {
            const data = await CreateChannel([ws.getUsername(), text])
            if (data.code !== 200) {
                switch (data.message) {
                    case ErrorUsernameDoesntExist:
                        errorMsg = "The inputted username does not exist."
                        break;
                }
                return invalidInput("create channel", errorMsg)
            }
            if (!data.data) {
                return invalidInput("create channel", errorMsg)
            }
            setChannels(value => [...value, data.data!.channel])
            ws.sendChannelAdded(ws.getUsername(), data.data.channel.id)
        } catch (err: unknown) {
            console.error(err)
            return invalidInput("create channel", errorMsg)
        }
    }

    return <Box flex={"1"} w={"100%"} pt={"10px"} overflow={"scroll"}>
        <ChannelsPageCreateChat onSubmit={onSubmit} />
        <Divider my={"18px"} />
        <Wrap spacing={"20px"}>
            {channels.map(r => <Channel key={r.id} channel={r} />)}
        </Wrap>
    </Box>
}

interface ChannelsPageCreateChatProps {
    onSubmit: (text: string) => Promise<void>
}

const ChannelsPageCreateChat = (props: ChannelsPageCreateChatProps) => {
    const [text, setText] = useState<string>("")

    return <HStack w={"100%"} justify={"center"}>
        <Input w={"50%"} value={text} onChange={(event) => setText(event.target.value)} />
        <Button onClick={() => {
            const sendText = text.trim()
            setText("")
            props.onSubmit(sendText)
        }}>Create Chat</Button>
    </HStack>
}