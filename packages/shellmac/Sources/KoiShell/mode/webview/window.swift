import AppKit
import SwiftUI
import WebView
import WebKit

struct KSWebViewOutput: Codable {}

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

class KSWebViewDelegate: NSObject, NSApplicationDelegate, NSWindowDelegate, WKScriptMessageHandler {
    var window: NSWindow!
    var webView: WKWebView!
    let name: String
    let url: String
    var appearance: NSAppearance!
    var initAppearance: NSAppearance!

    init(_ url: String, _ name: String) {
        self.url = url
        self.name = name
    }

    func applicationDidFinishLaunching(_ notification: Notification) {
        let enhanceURL = Bundle.module.url(forResource: "userscript", withExtension: "js")!
        let enhanceData = try! Data(contentsOf: enhanceURL)
        let enhanceRaw = String(decoding: enhanceData, as: UTF8.self)
            .replacingOccurrences(of: "KOISHELL_RUNTIME_SUPPORTS", with: "['enhance']")
        let userScript = WKUserScript(
            source: enhanceRaw,
            injectionTime: .atDocumentEnd,
            forMainFrameOnly: true
        )
        let userContentController = WKUserContentController()
        userContentController.addUserScript(userScript)
        userContentController.add(self, name: "shellmacHandler")
        let configuration = WKWebViewConfiguration()
        configuration.userContentController = userContentController
        configuration.preferences.setValue(true, forKey: "developerExtrasEnabled")
        webView = WKWebView(
            frame: .zero,
            configuration: configuration
        )
        webView.setValue(false, forKey: "drawsBackground")
        webView.load(URLRequest(url: URL(string: self.url)!))
        let webViewStore = WebViewStore(
            webView: webView
        )

        let contentView = KSWebView(url, webViewStore)

        window = NSWindow(
            contentRect: NSRect(x: 0, y: 0, width: 1366, height: 768),
            styleMask: [.titled, .closable, .miniaturizable, .resizable, .fullSizeContentView],
            backing: .buffered,
            defer: false
        )
        appearance = window.appearance
        initAppearance = window.appearance

        window.title = "\(self.name) - Cordis"
        window.titlebarAppearsTransparent = true
        // window.titleVisibility = .hidden
        // window.appearance = NSAppearance(named: NSAppearance.Name.vibrantLight)
        // window.isMovableByWindowBackground = true

        setMenu()

        window.center()
        window.setFrameAutosaveName("KSWebView")
        window.contentView = NSHostingView(rootView: contentView)
        window.delegate = self

        window.makeKeyAndOrderFront(nil)
        NSApp.activate(ignoringOtherApps: true)
    }

    func setMenu() {
        let mainMenu = NSMenu(title: "Application")
        NSApp.mainMenu = mainMenu

        let appMenuItem = mainMenu.addItem(withTitle: "Application", action: nil, keyEquivalent: "")
        let appMenu = NSMenu(title: "Application")
        appMenuItem.submenu = appMenu

        let appMenuAbout = appMenu.addItem(withTitle: "About Cordis Console", action: #selector(NSApp.orderFrontStandardAboutPanel(_:)), keyEquivalent: "")
        appMenuAbout.target = NSApp
        appMenu.addItem(NSMenuItem.separator())
        let appMenuServices = appMenu.addItem(withTitle: "Services", action: nil, keyEquivalent: "")
        appMenuServices.submenu = NSMenu(title: "Services")
        NSApp.servicesMenu = appMenuServices.submenu
        appMenu.addItem(NSMenuItem.separator())
        let appMenuHide = appMenu.addItem(withTitle: "Hide", action: #selector(NSApp.hide), keyEquivalent: "h")
        appMenuHide.target = NSApp
        let appMenuHideOthers = appMenu.addItem(withTitle: "Hide Others", action: #selector(NSApp.hideOtherApplications), keyEquivalent: "h")
        appMenuHideOthers.target = NSApp
        let appMenuShowAll = appMenu.addItem(withTitle: "Show All", action: #selector(NSApp.unhideAllApplications), keyEquivalent: "")
        appMenuShowAll.target = NSApp
        appMenu.addItem(NSMenuItem.separator())
        let appMenuQuit = appMenu.addItem(withTitle: "Quit Cordis Console", action: #selector(NSApp.terminate), keyEquivalent: "q")
        appMenuQuit.target = NSApp

        let consoleMenuItem = mainMenu.addItem(withTitle: "Console", action: nil, keyEquivalent: "")
        let consoleMenu = NSMenu(title: "Console")
        consoleMenuItem.submenu = consoleMenu
        consoleMenu.addItem(withTitle: "Close", action: #selector(NSWindow.performClose), keyEquivalent: "w")

        let editMenuItem = mainMenu.addItem(withTitle: "Edit", action: nil, keyEquivalent: "")
        let editMenu = NSMenu(title: "Edit")
        populateEditMenu(editMenu)
        editMenuItem.submenu = editMenu

        let viewMenuItem = mainMenu.addItem(withTitle: "View", action: nil, keyEquivalent: "")
        let viewMenu = NSMenu(title: "View")
        viewMenuItem.submenu = viewMenu
        viewMenu.addItem(withTitle: "Reload", action: #selector(refresh), keyEquivalent: "r")

        let goMenuItem = mainMenu.addItem(withTitle: "Go", action: nil, keyEquivalent: "")
        let goMenu = NSMenu(title: "Go")
        goMenuItem.submenu = goMenu
        let goDashboardMenuItem = goMenu.addItem(withTitle: "Dashboard", action: #selector(goDashboard), keyEquivalent: "1")
        // goDashboardMenuItem.image = NSImage(named: NSImage.homeTemplateName)
        let goPluginsMenuItem = goMenu.addItem(withTitle: "Plugins", action: #selector(goPlugins), keyEquivalent: "2")
        // goPluginsMenuItem.image = NSImage(named: NSImage.homeTemplateName)
        let goMarketMenuItem = goMenu.addItem(withTitle: "Market", action: #selector(goMarket), keyEquivalent: "3")
        // goMarketMenuItem.image = NSImage(named: NSImage.homeTemplateName)
        let goDependenciesMenuItem = goMenu.addItem(withTitle: "Dependencies", action: #selector(goDependencies), keyEquivalent: "4")
        // goDependenciesMenuItem.image = NSImage(named: NSImage.homeTemplateName)
        let goSandboxMenuItem = goMenu.addItem(withTitle: "Sandbox", action: #selector(goSandbox), keyEquivalent: "5")
        // goSandboxMenuItem.image = NSImage(named: NSImage.homeTemplateName)
        let goLogsMenuItem = goMenu.addItem(withTitle: "Logs", action: #selector(goLogs), keyEquivalent: "6")
        // goLogsMenuItem.image = NSImage(named: NSImage.homeTemplateName)

        let windowMenuItem = mainMenu.addItem(withTitle: "Window", action: nil, keyEquivalent: "")
        let windowMenu = NSMenu(title: "Window")
        populateWindowMenu(windowMenu)
        windowMenuItem.submenu = windowMenu
        NSApp.windowsMenu = windowMenu

        let instanceMenuItem = mainMenu.addItem(withTitle: name, action: nil, keyEquivalent: "")
        let instanceMenu = NSMenu(title: name)
        instanceMenuItem.submenu = instanceMenu
        instanceMenu.addItem(withTitle: "URL: \(url)", action: nil, keyEquivalent: "")

        let helpMenuItem = mainMenu.addItem(withTitle: "Help", action: nil, keyEquivalent: "")
        let helpMenu = NSMenu(title: "Help")
        populateHelpMenu(self, helpMenu)
        helpMenuItem.submenu = helpMenu
        NSApp.helpMenu = helpMenu
    }

    @objc
    func openDocumentation() {
        NSWorkspace.shared.open(URL(string: "https://koishi.chat")!)
    }

    @objc
    func refresh() {
        webView.reload()
    }

    @objc
    func goDashboard() {
        webView.load(URLRequest(url: URL(string: "\(url)/")!))
    }

    @objc
    func goPlugins() {
        webView.load(URLRequest(url: URL(string: "\(url)/plugins")!))
    }

    @objc
    func goMarket() {
        webView.load(URLRequest(url: URL(string: "\(url)/market")!))
    }

    @objc
    func goDependencies() {
        webView.load(URLRequest(url: URL(string: "\(url)/dependencies")!))
    }

    @objc
    func goSandbox() {
        webView.load(URLRequest(url: URL(string: "\(url)/sandbox")!))
    }

    @objc
    func goLogs() {
        webView.load(URLRequest(url: URL(string: "\(url)/logs")!))
    }

    func applicationShouldTerminateAfterLastWindowClosed(_ sender: NSApplication) -> Bool {
        NSApplication.shared.terminate(self)
        return true
    }

    func windowWillClose(_ notification: Notification) {
        print(try! JSONEncoder().encode(KSWebViewOutput()).base64EncodedString())
    }

    func userContentController(_ userContentController: WKUserContentController, didReceive message: WKScriptMessage) {
        guard let msg = message.body as? String else {
            log("Failed to parse shellmacHandler message \(message.body)")
            return
        }

        switch msg {
        case "TL":
            appearance = NSAppearance(named: NSAppearance.Name.vibrantLight)
        case "TD":
            appearance = NSAppearance(named: NSAppearance.Name.vibrantDark)
        case "TR":
            appearance = initAppearance
        default:
            break
        }

        setAppearance()
    }

    func setAppearance() {
        window.appearance = appearance
        window.appearanceSource.appearance = appearance
        window.invalidateShadow()
    }

    func windowDidBecomeKey(_ notification: Notification) {
        setAppearance()
    }

    func windowDidResignKey(_ notification: Notification) {
        setAppearance()
    }
}

struct KSWebView: View {
    @StateObject var webViewStore: WebViewStore
    private var url: String

    init(_ url: String, _ webViewStore: WebViewStore) {
        self.url = url
        _webViewStore = StateObject(
            wrappedValue: webViewStore
        )
    }

    var body: some View {
        ZStack {
            ZStack {}
                .frame(maxWidth: .infinity, maxHeight: .infinity)
                .visualEffectBackground()
                .ignoresSafeArea()

            ZStack {
                WebView(webView: webViewStore.webView)
                    .visualEffectBackground()
            }
        }
    }
}
