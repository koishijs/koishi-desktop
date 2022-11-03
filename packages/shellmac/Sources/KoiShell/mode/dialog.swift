import AppKit

struct KSDialogOutput: Codable {
    var result: Int
}

func ksDialog(_ arg: [String: Any]) {
    let alert = NSAlert()
    if let style = arg["style"] as? String {
        switch style {
        case "info":
            alert.alertStyle = .informational
        case "warn":
            alert.alertStyle = .warning
        case "error":
            alert.alertStyle = .critical
        default:
            break
        }
    }
    if let text1 = arg["text1"] as? String {
        alert.messageText = text1
    }
    if let text2 = arg["text2"] as? String {
        alert.informativeText = text2
    }
    if let buttonCount = arg["buttonCount"] as? Int {
        if buttonCount >= 1 {
            if let button1Text = arg["button1Text"] as? String {
                alert.addButton(withTitle: button1Text)
            } else {
                alert.addButton(withTitle: "OK")
            }
        }

        if buttonCount >= 2 {
            if let button2Text = arg["button2Text"] as? String {
                alert.addButton(withTitle: button2Text)
            } else {
                alert.addButton(withTitle: "Cancel")
            }
        }

        if buttonCount >= 3 {
            if let button3Text = arg["button3Text"] as? String {
                alert.addButton(withTitle: button3Text)
            } else {
                alert.addButton(withTitle: "Don’t Save")
            }
        }
    }

    var result = 0

    switch alert.runModal() {
    case NSApplication.ModalResponse.alertFirstButtonReturn:
        result = 1
    case NSApplication.ModalResponse.alertSecondButtonReturn:
        result = 2
    case NSApplication.ModalResponse.alertThirdButtonReturn:
        result = 3
    default:
        log("Unknown ModalResponse.")
        exit(EXIT_FAILURE)
    }

    print(try! JSONEncoder().encode(KSDialogOutput(result: 1)).base64EncodedString())
}
