import { ElNotification } from "element-plus";

/**
 * 在任意页面处理「对方发起通话」：通知 + 跳进对应单聊/群聊。
 * 若已在 ContactChat 且正是该会话，则交界面内逻辑处理（避免重复通知）。
 */
export function handleIncomingCallSignal(ev, router) {
  if (!ev?.data || typeof ev.data !== "string") return;
  const text = ev.data;
  let msg;
  try {
    msg = JSON.parse(text);
  } catch {
    return;
  }
  if (msg.type !== 3 || !msg.av_data) return;
  let av;
  try {
    av = JSON.parse(msg.av_data);
  } catch {
    return;
  }
  if (av.messageId !== "PROXY" || av.type !== "start_call") return;

  const needChatId = msg.receive_id?.startsWith("G")
    ? msg.receive_id
    : msg.send_id;
  const r = router.currentRoute.value;
  if (r.name === "ContactChat" && r.params?.id === needChatId) {
    return;
  }

  ElNotification({
    title: "通话请求",
    message: `${msg.send_name || "对方"} 发起通话，正在进入会话…`,
    type: "warning",
    duration: 6000,
  });
  if (needChatId) {
    router
      .push({ name: "ContactChat", params: { id: needChatId } })
      .catch(() => {});
  }
}

/** 先跑 ContactChat 内逻辑（可 return 掉非本会话信令），再跑全局来电跳转兜底 */
export function composeSocketOnMessage(router, inner) {
  return (ev) => {
    if (typeof inner === "function") {
      inner(ev);
    }
    handleIncomingCallSignal(ev, router);
  };
}
