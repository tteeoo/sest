function test() {
    
}
browser.tabs.executeScript({file: "/content_scripts/sest.js"})
.then(test)