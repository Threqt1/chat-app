import { WrapItem, Card, CardBody, Flex, Avatar, Spacer, Text } from "@chakra-ui/react"
import { useNavigate } from "react-router-dom"
import { TChannel } from "../../types"

export interface ChannelProps {
    channel: TChannel
}

export const Channel = (props: ChannelProps) => {
    const navigate = useNavigate()

    const date = new Date(props.channel.lastActivity.timestamp)

    return <WrapItem userSelect={"none"} w="100%" onClick={() => navigate("/message/" + props.channel.id)}>
        <Card size={"md"} w={"100%"} variant={"filled"} borderRadius={"2xl"} borderWidth={"1px"} boxShadow={"xl"}>
            <CardBody>
                <Flex w="100%" direction={"row"} align={"center"} gap="15px">
                    <Avatar name={props.channel.members.map(r => r.slice(0, 1)).join(" ")} size={"md"} />
                    <Text fontSize={"xl"} as={"b"}>{props.channel.members.join(", ")}</Text>
                    <Text>&bull;</Text>
                    <Text fontSize={"lg"}>{props.channel.lastActivity.content}</Text>
                    <Spacer />
                    <Text fontSize={"sm"}>{`${date.toLocaleDateString("en-US", { weekday: "short", month: "short", hour12: true, hour: "2-digit", year: "numeric", day: "2-digit", minute: "2-digit" })}`}</Text>
                </Flex>
            </CardBody>
        </Card>
    </WrapItem>
}