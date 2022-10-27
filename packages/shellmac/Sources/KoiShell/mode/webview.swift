import AppKit
import SwiftUI
import WebView
import WebKit

func ksWebView(_ arg: [String: Any]) {
    guard let url = arg["url"] as? String else {
        log("Failed to parse url.")
        exit(EXIT_FAILURE)
    }

    guard let name = arg["name"] as? String else {
        log("Failed to parse name.")
        exit(EXIT_FAILURE)
    }

    let delegate = KSWebViewDelegate(url, name)
    NSApp.delegate = delegate
    NSApp.run()
}

class KSWebViewDelegate: NSObject, NSApplicationDelegate, NSWindowDelegate {
    var window: NSWindow!
    var hostingView: NSView?
    var contentView: KSWebView
    let name: String

    init(_ url: String, _ name: String) {
        self.contentView = KSWebView(url)
        self.name = name
    }

    func applicationDidFinishLaunching(_ notification: Notification) {
        window = NSWindow(
            contentRect: NSRect(x: 0, y: 0, width: 1366, height: 768),
            styleMask: [.titled, .closable, .miniaturizable, .resizable, .fullSizeContentView],
            backing: .buffered,
            defer: false
        )

        window.title = "\(self.name) - Koishi"
        window.center()
        window.setFrameAutosaveName("KSWebView")
        hostingView = NSHostingView(rootView: contentView)
        window.contentView = hostingView
        window.delegate = self
        window.makeKeyAndOrderFront(nil)
        NSApp.activate(ignoringOtherApps: true)
    }

    func applicationShouldTerminateAfterLastWindowClosed(_ sender: NSApplication) -> Bool {
        NSApplication.shared.terminate(self)
        return true
    }
}

struct KSWebView: View {
    @StateObject var webViewStore: WebViewStore
    var url: String

    init(_ url: String) {
        self.url = url

        let enhanceURL = Bundle.module.url(forResource: "Resources/enhance", withExtension: "js")!
        let enhanceData = try! Data(contentsOf: enhanceURL)
        let enhanceRaw = String(decoding: enhanceData, as: UTF8.self)
        let userScript = WKUserScript(
            source: enhanceRaw,
            injectionTime: .atDocumentEnd,
            forMainFrameOnly: true
        )
        let userContentController = WKUserContentController()
        userContentController.addUserScript(userScript)
        let configuration = WKWebViewConfiguration()
        configuration.userContentController = userContentController
        _webViewStore = StateObject(
            wrappedValue: WebViewStore(
                webView: WKWebView(
                    frame: .zero,
                    configuration: configuration
                )
            )
        )
    }

    var body: some View {
        WebView(webView: webViewStore.webView)
            .onAppear {
                self.webViewStore.webView.setValue(false, forKey: "drawsBackground")
                self.webViewStore.webView.load(URLRequest(url: URL(string: self.url)!))
            }
    }
}
