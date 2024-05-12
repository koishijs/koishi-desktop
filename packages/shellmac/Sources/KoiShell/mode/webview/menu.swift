import Cocoa

@objc protocol EditMenuActions {
    func redo(_ sender: AnyObject)
    func undo(_ sender: AnyObject)
}

extension NSApplicationDelegate {
    func populateEditMenu(_ editMenu: NSMenu) {
        editMenu.addItem(withTitle: "Undo", action: #selector(EditMenuActions.undo), keyEquivalent: "z")
        editMenu.addItem(withTitle: "Redo", action: #selector(EditMenuActions.redo), keyEquivalent: "Z")
        editMenu.addItem(NSMenuItem.separator())
        editMenu.addItem(withTitle: "Cut", action: #selector(NSText.cut), keyEquivalent: "x")
        editMenu.addItem(withTitle: "Copy", action: #selector(NSText.copy), keyEquivalent: "c")
        editMenu.addItem(withTitle: "Paste", action: #selector(NSText.paste), keyEquivalent: "v")
        editMenu.addItem(withTitle: "Paste and Match Style", action: #selector(NSTextView.pasteAsRichText), keyEquivalent: "V")
        editMenu.addItem(withTitle: "Delete", action: #selector(NSText.delete), keyEquivalent: "\u{8}") // Backspace
        editMenu.addItem(withTitle: "Select All", action: #selector(NSText.selectAll), keyEquivalent: "a")
        editMenu.addItem(NSMenuItem.separator())

        let findMenuItem = editMenu.addItem(withTitle: "Find", action: nil, keyEquivalent: "")
        let findMenu = NSMenu(title: "Find")
        populateFindMenu(findMenu)
        findMenuItem.submenu = findMenu

        let spellingMenuItem = editMenu.addItem(withTitle: "Spelling", action: nil, keyEquivalent: "")
        let spellingMenu = NSMenu(title: "Spelling")
        populateSpellingMenu(spellingMenu)
        spellingMenuItem.submenu = spellingMenu
    }

    private func populateFindMenu(_ findMenu: NSMenu) {
        let findMenuItem = findMenu.addItem(withTitle: "Find…", action: #selector(NSResponder.performTextFinderAction), keyEquivalent: "f")
        findMenuItem.tag = NSTextFinder.Action.showFindInterface.rawValue

        let findNextMenuItem = findMenu.addItem(withTitle: "Find Next", action: #selector(NSResponder.performTextFinderAction), keyEquivalent: "g")
        findNextMenuItem.tag = NSTextFinder.Action.nextMatch.rawValue

        let findPreviousMenuItem = findMenu.addItem(withTitle: "Find Previous", action: #selector(NSResponder.performTextFinderAction), keyEquivalent: "G")
        findPreviousMenuItem.tag = NSTextFinder.Action.previousMatch.rawValue

        let findSelectionMenuItem = findMenu.addItem(withTitle: "Use Selection for Find", action: #selector(NSResponder.performTextFinderAction), keyEquivalent: "e")
        findSelectionMenuItem.tag = NSTextFinder.Action.setSearchString.rawValue

        findMenu.addItem(withTitle: "Jump to Selection", action: #selector(NSResponder.centerSelectionInVisibleArea), keyEquivalent: "j")
    }

    private func populateSpellingMenu(_ spellingMenu: NSMenu) {
        spellingMenu.addItem(withTitle: "Spelling…", action: #selector(NSText.showGuessPanel), keyEquivalent: ":")
        spellingMenu.addItem(withTitle: "Check Spelling", action: #selector(NSText.checkSpelling), keyEquivalent: ";")
        spellingMenu.addItem(withTitle: "Check Spelling as You Type", action: #selector(NSTextView.toggleContinuousSpellChecking), keyEquivalent: "")
    }

    func populateWindowMenu(_ windowMenu: NSMenu) {
        windowMenu.addItem(withTitle: "Minimize", action: #selector(NSWindow.performMiniaturize), keyEquivalent: "m")
        windowMenu.addItem(withTitle: "Zoom", action: #selector(NSWindow.performZoom), keyEquivalent: "")

        windowMenu.addItem(NSMenuItem.separator())

        let fullScreenMenuItem = windowMenu.addItem(withTitle: "Enter Full Screen", action: #selector(NSWindow.toggleFullScreen), keyEquivalent: "f")
        fullScreenMenuItem.keyEquivalentModifierMask = [.command, .control]

        windowMenu.addItem(NSMenuItem.separator())

        let allToFrontMenuItem = windowMenu.addItem(withTitle: "Bring All to Front", action: #selector(NSApp.arrangeInFront), keyEquivalent: "")
        allToFrontMenuItem.target = NSApp
    }

    func populateHelpMenu(_ del: KSWebViewDelegate, _ helpMenu: NSMenu) {
        let koishiDocumentationMenuItem = helpMenu.addItem(withTitle: "Cordis Documentation", action: #selector(del.openDocumentation), keyEquivalent: "h")
        koishiDocumentationMenuItem.keyEquivalentModifierMask = [.command, .shift]
    }
}
