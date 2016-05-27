const {dialog, app, Menu, Tray, shell, MenuItem} = require('electron');
const child_process = require('child_process')
const path = require('path')
const fs = require('fs')
const rimraf = require('rimraf')
const storage = require('electron-json-storage')

const powerSaveBlocker = require('electron').powerSaveBlocker;
powerSaveBlocker.start('prevent-app-suspension');

let appIcon = null;
let outPath = path.join(app.getPath('temp'), "assetkit")
console.log(outPath)
app.on('ready', () => {
  tray = appIcon = new Tray(path.join(app.getAppPath(), 'icon.png'));
  /*
  const contextMenu = Menu.buildFromTemplate([
      {label: "Build Project", click() { openProject() }},
      {type: 'separator'},
      {label: "Quit", click() { app.quit() }}
  ]);
  */
  appIcon.setToolTip('Asset Kit');
  //appIcon.setContextMenu(contextMenu);
  loadMenu()
});
//app.dock.hide()

function openProject() {
    selectedPath = dialog.showOpenDialog({properties: ["openDirectory"]})
    if (selectedPath) {
        saveOpenedProject(selectedPath[0], () => {
            buildProject(selectedPath[0])
        })
    }
}

function buildProject(project_path) {
    rimraf(outPath, () => {
        child_process.execFile(path.join(app.getAppPath(),"assetkit"), ["compile", project_path, outPath], (err, stdout, stderr) => {
            if (err) { throw err }
            console.log(stdout)
            console.log(stderr)
            loadMenu(() => {
                //dialog.showMessageBox({message: "編譯完成", buttons: ["OK"]})
                shell.openItem(outPath + "/index.html")
                //shell.showItemInFolder(outPath + "/index.html")
            })
        })
    })
}

function saveOpenedProject(path, cb) {
    storage.get('opened_projects', (err, data) => {
        if (err) { throw err }

        if (data["projects"] == undefined) {
            data["projects"] = []
        }

        if (!data["projects"].includes(path)) {
            data["projects"].push(path)
        }

        storage.set('opened_projects', data, cb)
    })
}

function loadMenu(cb) {
    const menu = new Menu()
    menu.append(new MenuItem({label:'Open Project', click: openProject}))
    menu.append(new MenuItem({type: 'separator'}))
    storage.get('opened_projects', (err, data) => {
        if (err) { throw err }

        if (data["projects"] != undefined) {
            data["projects"].forEach((p) => {
                menu.append(new MenuItem({label: p, submenu: [
                    //{label: "Build & Publish"},
                    {label: "Build", click() { buildProject(p) }}
                ]}))
            })
        }

        menu.append(new MenuItem({type: 'separator'}))
        menu.append(new MenuItem({label: "Quit", click() { app.quit() }}))

        appIcon.setContextMenu(menu)
        if (cb != undefined) {
            cb()
        }
    })
}
