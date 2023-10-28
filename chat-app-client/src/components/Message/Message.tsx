import { Card, CardBody, Flex, Avatar, Text, Box } from "@chakra-ui/react"

export type MessageProps = {
    avatarURL?: string
    name: string
    content: string
    pending?: boolean
    timestamp: number
}

export const Message = (props: MessageProps) => {
    return <Box w={"98%"}>
        <Card size={"sm"} w={"100%"} variant={"filled"} borderRadius={"2xl"} boxShadow={"xl"}>
            <CardBody>
                <Flex direction={"row"} gap={"15px"}>
                    <Avatar name={props.name} src={props.avatarURL} size={"md"} />
                    <Flex width={"100%"} direction={"column"}>
                        <Text fontSize={"xl"} as={"b"} opacity={props.pending ? 0.4 : 1}>{props.name}</Text>
                        <Text fontSize={"xl"} opacity={props.pending ? 0.4 : 1}>{props.content}</Text>
                    </Flex>
                </Flex>
            </CardBody>
        </Card>
    </Box>
}