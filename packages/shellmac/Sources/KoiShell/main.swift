import AppKit

signal(SIGINT, { _ in
    log("Received SIGINT. Shutting down...")
    exit(EXIT_SUCCESS)
})

NSApplication.shared.setActivationPolicy(.regular)

if CommandLine.arguments.count != 2 {
    log("argc not valid.")
    exit(EXIT_FAILURE)
}

let rawEncodedArg = CommandLine.arguments[1]
guard let rawArg = Data(base64Encoded: rawEncodedArg) else {
    log("Failed to parse encoded arg.")
    exit(EXIT_FAILURE)
}

guard let argObj = try? JSONSerialization.jsonObject(with: Data(rawArg)) else {
    log("Failed to parse arg.")
    exit(EXIT_FAILURE)
}
guard let arg = argObj as? [String: Any] else {
    log("Failed to parse arg.")
    exit(EXIT_FAILURE)
}
guard let mode = arg["mode"] as? String else {
    log("Failed to parse mode.")
    exit(EXIT_FAILURE)
}

switch mode {
case "webview":
    ksWebView(arg)
default:
    log("Unknown mode: \(mode)")
    exit(EXIT_FAILURE)
}
