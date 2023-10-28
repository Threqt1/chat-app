import { Avatar, Card, CardBody, Divider, Flex, HStack, Icon, Text, ToastId, useToast } from "@chakra-ui/react"
import { useEffect, useRef } from "react"
import { RiContactsBookLine, RiSettings5Line } from "react-icons/ri"
import { Outlet, useLoaderData, useNavigate } from "react-router-dom"
import { TMessage } from "./types";
import { WebsocketLoaderResult } from "./main";

const MAX_TOASTS = 5;
const TRUNC_NAME_LENGTH = 10;
const TRUNC_MESSAGE_LENGTH = 15;

export const MainPage = () => {
    const { ws } = useLoaderData() as WebsocketLoaderResult

    const navigate = useNavigate()
    const notificationToast = useToast();
    const toastsArray = useRef<ToastId[]>([])

    useEffect(() => {
        const notificationCallback = (chat: TMessage) => {
            if (toastsArray.current.length >= MAX_TOASTS) {
                notificationToast.close(toastsArray.current[0]);
                toastsArray.current = toastsArray.current.slice(1);
            }
            toastsArray.current = [...toastsArray.current, notificationToast({
                position: 'top-right',
                render: () => (
                    <Card size={"md"} w={"100%"} variant={"filled"} borderRadius={"2xl"} borderWidth={"1px"} boxShadow={"xl"}>
                        <CardBody>
                            <Flex w="100%" direction={"row"} align={"center"} gap="15px">
                                <Avatar name={"krish parikh"} size={"md"} />
                                <Text fontSize={"xl"} as={"b"}>{chat.from.length > TRUNC_NAME_LENGTH ? chat.from.slice(0, TRUNC_NAME_LENGTH) + "..." : chat.from}</Text>
                                <Text>&bull;</Text>
                                <Text fontSize={"lg"}>{chat.content.length > TRUNC_MESSAGE_LENGTH ? chat.content.slice(0, TRUNC_MESSAGE_LENGTH) + "..." : chat.content}</Text>
                            </Flex>
                        </CardBody>
                    </Card>
                ),
            })]
        }

        ws.getEmitter().on("chat_notification", notificationCallback)

        return () => {
            ws.getEmitter().off("chat_notification", notificationCallback)
        }
    }, [notificationToast, ws])

    return (
        <Flex h={"100vh"} w={"100vw"} p={"5px"} direction={"column"} align={"center"}>
            <HStack divider={<Text px={"10px"} >&bull;</Text>} align={"center"} userSelect={"none"}>
                <Icon my={"10px"} as={RiContactsBookLine} boxSize={10} onClick={() => navigate("/channels")} />
                <Icon my={"10px"} as={RiSettings5Line} boxSize={10} onClick={() => navigate("/settings")} />
            </HStack>
            <Divider my={"10px"} />
            <Outlet />
        </Flex >
    )
}