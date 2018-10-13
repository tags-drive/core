// For every fetch request we have to use `credentials: "same-origin"` to pass Cookie

/* Global */
const sortType = {
    name: "name",
    size: "size",
    time: "time"
};

const sortOrder = {
    asc: "asc",
    desc: "desc"
};

const logTypes = {
    info: "info",
    error: "error"
};

// GlobalStore contains global data
var GlobalStore = {
    data: {
        allFiles: [],
        allTags: []
    },
    updateFiles: function() {
        fetch("/api/files", {
            method: "GET",
            credentials: "same-origin"
        })
            .then(data => data.json())
            .then(files => this.setFiles(files))
            .catch(err => console.error(err)); // user can live without this error, so we won't use logError() here
    },
    updateTags: function() {
        fetch("/api/tags", {
            method: "GET",
            credentials: "same-origin"
        })
            .then(data => data.json())
            .then(tags => {
                this.data.allTags = tags;
            })
            .catch(err => console.error(err)); // user can live without this error, so we won't use logError() here
    },
    setFiles: function(files) {
        // Change time from "2018-08-23T22:48:59.0459184+03:00" to "23-08-2018 22:48"
        for (let i in files) {
            files[i].addTime = new Date(files[i].addTime).format("dd-mm-yyyy HH:MM");
        }
        this.data.allFiles = files;
    }
};

var GlobalState = {
    mainBlockOpacity: 1,
    showDropLayer: true, // when we show modal-window with tags showDropLayer is false
    selectMode: false
};

// Init should be called onload
function Init() {
    GlobalStore.updateTags();
    GlobalStore.updateFiles();
}

function updateStore() {
    GlobalStore.updateFiles();
    GlobalStore.updateTags();
}

function isErrorStatusCode(statusCode) {
    if (400 <= statusCode && statusCode < 600) {
        return true;
    }
    return false;
}

function logInfo(msg) {
    console.log(msg);
    logWindow.add(logTypes.info, msg);
}

function logError(err) {
    console.error(err);
    logWindow.add(logTypes.error, err);
}

/* Main instances */

// Top bar
var topBar = new Vue({
    el: "#top-bar",
    mixins: [VueClickaway.mixin],
    data: {
        // Tag search
        tagPrefix: "",
        showTagsList: false,
        pickedTags: [],
        unusedTags: [],
        // Advanced search
        text: "",
        selectedMode: "And"
    },
    methods: {
        tagsMenu: function() {
            return {
                show: () => {
                    if (this.pickedTags.length == 0) {
                        // Need to fill unusedTags
                        this.unusedTags = [];
                        for (let tag in GlobalStore.data.allTags) {
                            this.unusedTags.push(GlobalStore.data.allTags[tag]);
                        }
                    }

                    this.showTagsList = true;
                },
                hide: () => {
                    this.showTagsList = false;
                }
            };
        },
        search: function() {
            return {
                usual: () => {
                    let params = new URLSearchParams();
                    // tags
                    if (this.pickedTags.length != 0) {
                        let tags = [];
                        for (let tag of this.pickedTags) {
                            tags.push(tag.id);
                        }
                        params.append("tags", tags.join(","));
                    }
                    // search
                    if (this.text != "") {
                        params.append("search", this.text);
                    }
                    // mode
                    params.append("mode", this.selectedMode.toLowerCase());
                    // Can skip sort and order, because server will use default values

                    fetch("/api/files?" + params, {
                        method: "GET",
                        credentials: "same-origin"
                    })
                        .then(data => data.json())
                        .then(files => {
                            GlobalStore.setFiles(files);
                            // Reset sortParams
                            mainBlock.sort().reset();
                        });
                },
                advanced: (sType, sOrder) => {
                    let params = new URLSearchParams();
                    // tags
                    if (this.pickedTags.length != 0) {
                        let tags = [];
                        for (let tag of this.pickedTags) {
                            tags.push(tag.name);
                        }
                        params.append("tags", tags.join(","));
                    }
                    // search
                    if (this.text != "") {
                        params.append("search", this.text);
                    }
                    // sort
                    params.append("sort", sType);
                    // order
                    params.append("order", sOrder);
                    // mode
                    params.append("mode", this.selectedMode.toLowerCase());

                    fetch("/api/files?" + params, {
                        method: "GET",
                        credentials: "same-origin"
                    })
                        .then(data => data.json())
                        .then(files => GlobalStore.setFiles(files));
                }
            };
        },
        input: function() {
            return {
                tags: {
                    add: tagID => {
                        let index = -1;
                        for (let i in this.unusedTags) {
                            if (this.unusedTags[i].id == tagID) {
                                index = i;
                                break;
                            }
                        }
                        if (index == -1) {
                            return;
                        }

                        // Add a tag into pickedTags
                        this.pickedTags.push(this.unusedTags[index]);
                        // Remove a tag
                        this.unusedTags.splice(index, 1);
                    },
                    delete: tagID => {
                        let index = -1;
                        for (i in this.pickedTags) {
                            if (this.pickedTags[i].id == tagID) {
                                index = i;
                                break;
                            }
                        }
                        if (index == -1) {
                            return;
                        }

                        // Return a tag to unusedTags
                        this.unusedTags.push(this.pickedTags[index]);
                        // Remove an element
                        this.pickedTags.splice(index, 1);
                    }
                }
            };
        },
        settings: function() {
            return {
                tags: () => modalWindow.showWindow().globalTagsUpdating(),
                logout: () => {
                    if (!confirm("Are you sure you want log out?")) {
                        return;
                    }

                    fetch("/logout", {
                        method: "POST",
                        credentials: "same-origin"
                    })
                        .then(resp => {
                            if (isErrorStatusCode(resp.status)) {
                                resp.text().then(text => {
                                    logError(text);
                                });
                                return;
                            }

                            location.reload(true);
                        })
                        .catch(err => logError(err));
                }
            };
        }
    }
});

// Main block
var mainBlock = new Vue({
    el: "#main-block",
    data: {
        sharedData: GlobalStore.data,
        sharedState: GlobalState,

        sortByNameModeAsc: true,
        sortBySizeModeAsc: true,
        sortByTimeModeAsc: true,
        lastSortType: sortType.name,

        allSelected: false,
        selectCount: 0
    },
    methods: {
        // Sorts
        sort: function() {
            return {
                byName: () => {
                    if (this.lastSortType == sortType.name) {
                        this.sortByNameModeAsc = !this.sortByNameModeAsc;
                    } else {
                        // Use default settings
                        this.sort().reset();
                    }
                    this.lastSortType = sortType.name;

                    let type = sortType.name,
                        order = this.sortByNameModeAsc ? sortOrder.asc : sortOrder.desc;

                    topBar.search().advanced(type, order);
                },
                bySize: () => {
                    if (this.lastSortType == sortType.size) {
                        this.sortBySizeModeAsc = !this.sortBySizeModeAsc;
                    } else {
                        // Use default settings
                        this.sort().reset();
                    }
                    this.lastSortType = sortType.size;

                    let type = sortType.size,
                        order = this.sortBySizeModeAsc ? sortOrder.asc : sortOrder.desc;

                    topBar.search().advanced(type, order);
                },
                byTime: () => {
                    if (this.lastSortType == sortType.time) {
                        this.sortByTimeModeAsc = !this.sortByTimeModeAsc;
                    } else {
                        // Use default settings
                        this.sort().reset();
                    }
                    this.lastSortType = sortType.time;

                    let type = sortType.time,
                        order = this.sortByTimeModeAsc ? sortOrder.asc : sortOrder.desc;

                    topBar.search().advanced(type, order);
                },
                reset: () => {
                    this.sortByNameModeAsc = true;
                    this.sortBySizeModeAsc = true;
                    this.sortByTimeModeAsc = true;
                }
            };
        },
        // Select mode
        toggleAllFiles: function() {
            if (!this.allSelected) {
                this.selectCount = this.sharedData.allFiles.length;
                this.allSelected = true;
                GlobalState.selectMode = true;

                for (i in this.$children) {
                    this.$children[i].select();
                }
            } else {
                this.selectCount = 0;
                this.allSelected = false;
                GlobalState.selectMode = false;

                for (i in this.$children) {
                    this.$children[i].unselect();
                }
            }
        },
        unselectAllFile: function() {
            for (i in this.$children) {
                this.$children[i].unselect();
            }
            this.allSelected = false;
            GlobalState.selectMode = false;
            this.selectCount = 0;
        },
        getSelectedFiles: function() {
            let files = [];
            for (i in this.$children) {
                if (this.$children[i].selected) {
                    files.push(this.$children[i].file);
                }
            }
            return files;
        },
        // For children
        selectFile: function() {
            this.selectCount++;
            GlobalState.selectMode = true;
            if (this.selectCount == this.sharedData.allFiles.length) {
                this.allSelected = true;
            }
        },
        unselectFile: function() {
            this.selectCount--;
            this.allSelected = false;
            if (this.selectCount == 0) {
                GlobalState.selectMode = false;
            }
        }
    },
    template: `
	<table :style="{'opacity': sharedState.mainBlockOpacity}" class="file-table" style="width:100%;">
		<tr style="top: 100px;">
			<th style="text-align: center; width: 30px;">
				<input
				type="checkbox"
				:indeterminate.prop="selectCount > 0 && selectCount != sharedData.allFiles.length"
				v-model="allSelected"
				@click="toggleAllFiles"
				style="height: 15px; width: 15px;"
				title="Select all"
			>
			</th>
			<!-- Preview image -->
			<th></th> 
			<th>
				Filename
				<i class="material-icons" id="sortByNameIcon" @click="sort().byName()" :style="[sortByNameModeAsc ? {'transform': 'scale(1, 1)'} : {'transform': 'scale(1, -1)'}]" style="font-size: 20px; cursor: pointer;">
					sort
				</i>
			</th>
			<th>Tags</th>
			<th>
				Size (MB)
				<i class="material-icons" id="sortByNameSize" @click="sort().bySize()" :style="[sortBySizeModeAsc ? {'transform': 'scale(1, 1)'} : {'transform': 'scale(1, -1)'}]" style="transform: scale(1, 1); font-size: 20px; cursor: pointer;">
					sort
				</i>
			</th>
			<th>
				Time of adding
				<i class="material-icons" id="sortByNameTime" @click="sort().byTime()" :style="[sortByTimeModeAsc ? {'transform': 'scale(1, 1)'} : {'transform': 'scale(1, -1)'}]" style="transform: scale(1, 1); font-size: 20px; cursor: pointer;">
					sort
				</i>
			</th>
		</tr>
		<files v-for="file in sharedData.allFiles" :file="file" :allTags="sharedData.allTags"></files>
	</table>`
});

/* Secondary instances */

// Upload block
var uploader = new Vue({
    el: "#upload-block",
    data: {
        sharedData: GlobalStore.data,
        counter: 0 // for definition did user drag file into div. If counter > 0, user dragged file.
    },
    created() {
        // Add listeners
        document.ondragenter = () => {
            if (GlobalState.showDropLayer) {
                this.counter++;
            }
        };
        document.ondragleave = () => {
            if (GlobalState.showDropLayer) {
                this.counter--;
            }
        };
        document.ondrop = () => {
            if (GlobalState.showDropLayer) {
                this.counter = 0;
            }
        };
        setInterval(() => {
            if (this.counter == 0) {
                GlobalState.mainBlockOpacity = 1;
            } else {
                GlobalState.mainBlockOpacity = 0.3;
            }
        }, 10);
    },
    methods: {
        upload: function(event) {
            var formData = new FormData();

            for (file of event.dataTransfer.files) {
                formData.append("files", file, file.name);
            }

            fetch("/api/files", {
                body: formData,
                method: "POST",
                credentials: "same-origin"
            })
                .then(resp => {
                    if (isErrorStatusCode(resp.status)) {
                        resp.text().then(text => {
                            logError(text);
                        });
                        return;
                    }
                    // Update list of files
                    updateStore();
                    return resp.json();
                })
                .then(log => {
                    if (log === undefined) {
                        return;
                    }
                    console.log(log);
                    /* Schema:
                    [
                        {
                            filename: string,
                            isError: boolean,
                            error: string (when isError == true),
                            status: string (when isError == false)
                        }
                    ]
                    */
                    for (let i in log) {
                        let msg = log[i].filename;
                        if (log[i].isError) {
                            msg += " " + log[i].error;
                        } else {
                            msg += " " + log[i].status;
                        }

                        if (log[i].isError) {
                            logError(msg);
                        } else {
                            logInfo(msg);
                        }
                    }
                })
                .catch(err => logError(err));
        }
    }
});

// Context menu (right click on a file)
var contextMenu = new Vue({
    el: "#context-menu",
    mixins: [VueClickaway.mixin], // from vue-clickaway
    data: {
        sharedState: GlobalState,
        file: null,
        // Style
        top: "0px",
        left: "0px",
        show: false,
        // File changing
        newName: "",
        newTags: [],
        description: "",
        // For calculation of position
        divWidth: 140,
        divHeight: 125
    },
    methods: {
        setFile: function(file) {
            this.file = file;
        },
        // UI
        showMenu: function(x, y) {
            const offset = 10;
            x += offset;
            y += offset;
            if (x + this.divWidth > window.innerWidth) {
                x -= offset * 2;
                x -= this.divWidth;
            }
            if (y + this.divHeight > window.innerHeight) {
                y -= offset * 2;
                y -= this.divHeight;
            }
            this.left = x + "px";
            this.top = y + "px";
            this.show = true;
        },
        hideMenu: function() {
            this.show = false;
        },
        // Options of context menu
        regularMode: function() {
            return {
                changeName: () => {
                    this.show = false;
                    modalWindow.showWindow().regularRenaming(this.file);
                },
                changeTags: () => {
                    this.show = false;
                    modalWindow.showWindow().regularFileTagsUpdating(this.file);
                },
                changeDescription: () => {
                    this.show = false;
                    modalWindow.showWindow().regularDescriptionChanging(this.file);
                },
                deleteFile: () => {
                    this.show = false;
                    modalWindow.showWindow().regularDeleting(this.file);
                }
            };
        },
        // Options of context menu (select mode)
        selectMode: function() {
            return {
                addTags: () => {
                    this.show = false;
                    modalWindow.showWindow().selectFilesTagsAdding(mainBlock.getSelectedFiles());
                },
                deleteTags: () => {
                    this.show = false;
                    modalWindow.showWindow().selectFilesTagsDeleting(mainBlock.getSelectedFiles());
                },
                downloadFiles: () => {
                    let params = new URLSearchParams();
                    let files = mainBlock.getSelectedFiles();
                    for (let file of files) {
                        // need to use link to a file, not filename
                        params.append("file", file.origin);
                    }

                    fetch("/api/files/download?" + params, {
                        method: "GET",
                        credentials: "same-origin"
                    }).then(resp => {
                        resp.blob().then(file => {
                            let a = document.createElement("a"),
                                url = URL.createObjectURL(file);

                            a.href = url;
                            a.download = "files.zip";
                            document.body.appendChild(a);
                            a.click();

                            document.body.removeChild(a);
                            window.URL.revokeObjectURL(url);
                        });
                    });
                },
                deleteFiles: () => {
                    this.show = false;
                    modalWindow.showWindow().selectDeleting(mainBlock.getSelectedFiles());
                }
            };
        }
    }
});

// Modal window
// It's called from a context menu
var modalWindow = new Vue({
    el: "#modal-window",
    data: {
        file: null,
        show: false,
        selectedFiles: [],
        sharedData: GlobalStore.data,
        // Modes
        globalTagsMode: false,
        //
        regularRenameMode: false,
        regularFileTagsMode: false,
        regularDescriptionMode: false,
        regularDeleteMode: false,
        //
        selectFilesTagsAddMode: false,
        selectFilesTagsDeleteMode: false,
        selectDeleteMode: false,
        // For files API
        fileNewData: {
            newFilename: "",
            unusedTags: [],
            newTags: [],
            newDescription: ""
        },
        // For tags API
        newTag: {}
    },
    methods: {
        // UI
        showWindow: function() {
            return {
                globalTagsUpdating: () => {
                    GlobalState.showDropLayer = false;

                    this.globalTagsMode = true;

                    this.show = true;
                },
                // Regular mode
                regularRenaming: file => {
                    GlobalState.showDropLayer = false;

                    this.file = file;
                    this.regularRenameMode = true;
                    this.fileNewData.newFilename = file.filename;

                    this.show = true;
                },
                regularFileTagsUpdating: file => {
                    GlobalState.showDropLayer = false;

                    this.fileNewData.newTags = [];
                    this.fileNewData.unusedTags = [];

                    for (let id in this.sharedData.allTags) {
                        if (file.tags.includes(Number(id))) {
                            this.fileNewData.newTags.push(this.sharedData.allTags[id]);
                        } else {
                            this.fileNewData.unusedTags.push(this.sharedData.allTags[id]);
                        }
                    }

                    this.file = file;
                    this.regularFileTagsMode = true;

                    this.show = true;
                },
                regularDescriptionChanging: file => {
                    GlobalState.showDropLayer = false;

                    this.file = file;
                    this.fileNewData.unusedTags;
                    this.regularDescriptionMode = true;

                    this.show = true;
                },
                regularDeleting: file => {
                    GlobalState.showDropLayer = false;

                    this.file = file;
                    this.regularDeleteMode = true;

                    this.show = true;
                },
                // Select mode
                selectFilesTagsAdding: files => {
                    GlobalState.showDropLayer = false;

                    this.selectedFiles = files;
                    this.selectFilesTagsAddMode = true;

                    this.show = true;
                },
                selectFilesTagsDeleting: files => {
                    GlobalState.showDropLayer = false;

                    this.selectedFiles = files;
                    this.selectFilesTagsDeleteMode = true;

                    this.show = true;
                },
                selectDeleting: files => {
                    GlobalState.showDropLayer = false;

                    this.selectedFiles = files;
                    this.selectDeleteMode = true;

                    this.show = true;
                }
            };
        },
        hideWindow: function() {
            this.globalTagsMode = false;
            this.regularRenameMode = false;
            this.regularFileTagsMode = false;
            this.regularDescriptionMode = false;
            this.regularDeleteMode = false;
            this.selectFilesTagsAddMode = false;
            this.selectFilesTagsDeleteMode = false;
            this.selectDeleteMode = false;
            this.show = false;

            GlobalState.showDropLayer = true;
        },
        // Drag and drop
        tagsDragAndDrop: function() {
            return {
                add: ev => {
                    let tagID = Number(ev.dataTransfer.getData("tagName"));
                    let index = -1;
                    for (i in this.fileNewData.unusedTags) {
                        if (this.fileNewData.unusedTags[i].id == tagID) {
                            index = i;
                            break;
                        }
                    }
                    if (index == -1) {
                        return;
                    }
                    this.fileNewData.newTags.push(this.fileNewData.unusedTags[index]);
                    this.fileNewData.unusedTags.splice(index, 1);
                },
                del: ev => {
                    let tagID = ev.dataTransfer.getData("tagName");
                    let index = -1;
                    for (i in this.fileNewData.newTags) {
                        if (this.fileNewData.newTags[i].id == tagID) {
                            index = i;
                            break;
                        }
                    }
                    if (index == -1) {
                        return;
                    }
                    this.fileNewData.unusedTags.push(this.fileNewData.newTags[index]);
                    this.fileNewData.newTags.splice(index, 1);
                }
            };
        },
        // Files API
        filesAPI: function() {
            return {
                // Regular mode
                rename: () => {
                    let params = new URLSearchParams();
                    params.append("file", this.file.filename);
                    params.append("new-name", this.fileNewData.newFilename);

                    fetch("/api/files", {
                        method: "PUT",
                        body: params,
                        credentials: "same-origin"
                    })
                        .then(resp => {
                            if (isErrorStatusCode(resp.status)) {
                                resp.text().then(text => {
                                    logError(text);
                                });
                                return;
                            }
                            // Refresh list of files
                            topBar.search().usual();
                            this.hideWindow();
                        })
                        .catch(err => {
                            logError(err);
                        });
                },
                updateTags: () => {
                    let params = new URLSearchParams();
                    let tags = this.fileNewData.newTags.map(tag => tag.id);
                    params.append("file", this.file.filename);
                    let tagsParam = "empty";
                    if (tags.length != 0) {
                        tagsParam = tags.join(",");
                    }
                    params.append("tags", tagsParam);

                    fetch("/api/files", {
                        method: "PUT",
                        body: params,
                        credentials: "same-origin"
                    })
                        .then(resp => {
                            if (isErrorStatusCode(resp.status)) {
                                resp.text().then(text => {
                                    logError(text);
                                });
                                return;
                            }
                            // Refresh list of files
                            topBar.search().usual();
                            this.hideWindow();
                        })
                        .catch(err => {
                            logError(err);
                        });
                },
                updateDescription: () => {
                    let params = new URLSearchParams();
                    params.append("file", this.file.filename);
                    params.append("description", this.fileNewData.newDescription);

                    fetch("/api/files", {
                        method: "PUT",
                        body: params,
                        credentials: "same-origin"
                    })
                        .then(resp => {
                            if (isErrorStatusCode(resp.status)) {
                                resp.text().then(text => {
                                    logError(text);
                                });
                                return;
                            }
                            // Refresh list of files
                            topBar.search().usual();
                            this.hideWindow();
                        })
                        .catch(err => {
                            logError(err);
                        });
                },
                deleteFile: () => {
                    let params = new URLSearchParams();
                    params.append("file", this.file.filename);

                    fetch("/api/files?" + params, {
                        method: "DELETE",
                        credentials: "same-origin"
                    })
                        .then(resp => {
                            if (isErrorStatusCode(resp.status)) {
                                resp.text().then(text => {
                                    logError(text);
                                });
                                return;
                            }

                            // Refresh list of files
                            updateStore();
                            this.hideWindow();
                            return resp.json();
                        })
                        .then(log => {
                            if (log === undefined) {
                                return;
                            }
                            console.log(log);
                            /* Schema:
                            [
                                {
                                    filename: string,
                                    isError: boolean,
                                    error: string (when isError == true),
                                    status: string (when isError == false)
                                }
                            ]
                            */
                            for (let i in log) {
                                let msg = log[i].filename;
                                if (log[i].isError) {
                                    msg += " " + log[i].error;
                                } else {
                                    msg += " " + log[i].status;
                                }

                                if (log[i].isError) {
                                    logError(msg);
                                } else {
                                    logInfo(msg);
                                }
                            }
                        })
                        .catch(err => logError(err));
                },
                // Select mode
                addSelectedFilesTags: tagIDs => {
                    if (tagIDs.length == 0) {
                        return;
                    }

                    // Update tags and refresh list of files after all changes
                    (async () => {
                        for (file of this.selectedFiles) {
                            let tags = new Set(file.tags);
                            for (tag of tagIDs) {
                                tags.add(tag);
                            }

                            let params = new URLSearchParams();
                            params.append("file", file.filename);
                            params.append("tags", Array.from(tags).join(","));

                            await fetch("/api/files", {
                                method: "PUT",
                                body: params
                            })
                                .then(resp => {
                                    if (isErrorStatusCode(resp.status)) {
                                        resp.text().then(text => {
                                            logError(text);
                                        });
                                        return;
                                    }
                                })
                                .catch(err => logError(err));
                        }
                    })()
                        .then(() => {
                            // Refresh list of files
                            mainBlock.unselectAllFile();
                            topBar.search().usual();
                            this.hideWindow();
                        })
                        .catch(err => logError(err));
                },
                deleteSelectedFilesTags: tagIDs => {
                    if (tagIDs.length == 0) {
                        return;
                    }

                    // Update tags and refresh list of files after all changes
                    (async () => {
                        for (file of this.selectedFiles) {
                            let tags = new Set(file.tags);
                            for (tag of tagIDs) {
                                tags.delete(tag);
                            }

                            let params = new URLSearchParams();
                            let tagsList = "empty";
                            if (tags.size != 0) {
                                tagsList = Array.from(tags).join(",");
                            }
                            params.append("file", file.filename);
                            params.append("tags", tagsList);

                            await fetch("/api/files", {
                                method: "PUT",
                                body: params
                            })
                                .then(resp => {
                                    if (isErrorStatusCode(resp.status)) {
                                        resp.text().then(text => {
                                            logError(text);
                                        });
                                    }
                                })
                                .catch(err => logError(err));
                        }
                    })()
                        .then(() => {
                            // Refresh list of files
                            mainBlock.unselectAllFile();
                            topBar.search().usual();
                            this.hideWindow();
                        })
                        .catch(err => logError(err));
                },
                deleteSelectedFiles: () => {
                    let params = new URLSearchParams();
                    for (f of this.selectedFiles) {
                        params.append("file", f.filename);
                    }

                    fetch("/api/files?" + params, {
                        method: "DELETE",
                        credentials: "same-origin"
                    })
                        .then(resp => {
                            if (isErrorStatusCode(resp.status)) {
                                resp.text().then(text => {
                                    logError(text);
                                });
                                return;
                            }

                            // Refresh list of files
                            updateStore();
                            this.hideWindow();
                            return resp.json();
                        })
                        .then(log => {
                            if (log === undefined) {
                                return;
                            }
                            console.log(log);
                            /* Schema:
                            [
                                {
                                    filename: string,
                                    isError: boolean,
                                    error: string (when isError == true),
                                    status: string (when isError == false)
                                }
                            ]
                            */
                            for (let i in log) {
                                let msg = log[i].filename;
                                if (log[i].isError) {
                                    msg += " " + log[i].error;
                                } else {
                                    msg += " " + log[i].status;
                                }

                                if (log[i].isError) {
                                    logError(msg);
                                } else {
                                    logInfo(msg);
                                }
                            }
                        })
                        .catch(err => logError(err));

                    // If we don't call this function, next files will become selected.
                    mainBlock.unselectAllFile();
                }
            };
        },
        // Tags API
        tagsAPI: function() {
            return {
                // Requests
                add: (name, color) => {
                    let params = new URLSearchParams();
                    params.append("name", name);
                    params.append("color", color);

                    fetch("/api/tags", {
                        method: "POST",
                        body: params,
                        credentials: "same-origin"
                    })
                        .then(resp => {
                            if (isErrorStatusCode(resp.status)) {
                                resp.text().then(text => {
                                    logError(text);
                                });
                                return;
                            }

                            this.tagsAPI().delNewTag();
                            GlobalStore.updateTags();
                        })
                        .catch(err => {
                            logError(err);
                        });
                },
                change: (tagID, newName, newColor) => {
                    let params = new URLSearchParams();
                    params.append("id", tagID);
                    params.append("name", newName);
                    params.append("color", newColor);

                    fetch("/api/tags", {
                        method: "PUT",
                        body: params,
                        credentials: "same-origin"
                    })
                        .then(resp => {
                            if (isErrorStatusCode(resp.status)) {
                                resp.text().then(text => {
                                    logError(text);
                                });
                                return;
                            }

                            GlobalStore.updateTags();
                        })
                        .catch(err => {
                            logError(err);
                        });
                },
                del: tagID => {
                    let params = new URLSearchParams();
                    params.append("id", tagID);

                    fetch("/api/tags?" + params, {
                        method: "DELETE",
                        credentials: "same-origin"
                    })
                        .then(resp => {
                            if (isErrorStatusCode(resp.status)) {
                                resp.text().then(text => {
                                    logError(text);
                                });
                                return;
                            }

                            GlobalStore.updateTags();
                            // Need to update files to remove deleted tag
                            topBar.search().usual();
                            return resp.text();
                        })
                        .catch(err => {
                            logError(err);
                        });
                },
                // delNewTag deletes tag from tagsNewData.newTag
                delNewTag: () => {
                    this.newTag = {};
                }
            };
        }
    }
});

var logWindow = new Vue({
    el: "#log-window",
    data: {
        sharedLogTypes: logTypes,
        // Const
        hideTimeout: 5 * 1000, // 5s in milliseconds
        // States
        show: false,
        isMouseInside: false, // if isMouseInside, hideAfter isn't changed
        hideAfter: 5 * 1000,
        // UI
        opacity: 1,
        lastScrollHeight: 0,
        // Data
        /* events - array of objects:
           {
             type: string,
             msg: string,
             time: string
           }
        */
        events: []
    },
    created: function() {
        const msTimeout = 50; // 50ms is good enough. When 100, FPS is too low

        setInterval(() => {
            if (this.isMouseInside) {
                return;
            }
            if (this.hideAfter < 0) {
                this.show = false;
            }
            this.hideAfter -= msTimeout;
            if (this.hideAfter < 1000) {
                this.opacity = this.hideAfter / 1000;
            }
        }, msTimeout);
    },
    methods: {
        // UI
        window: function() {
            return {
                show: () => {
                    this.opacity = 1;
                    this.hideAfter = this.hideTimeout;
                    this.show = true;
                },
                hide: () => {
                    this.show = false;
                    this.isMouseInside = false;
                },
                mouseEnter: () => {
                    this.isMouseInside = true;
                    this.window().show(); // update opacity and hideAfter
                },
                mouseLeave: () => {
                    this.isMouseInside = false;
                },
                scrollToEnd: () => {
                    this.$el.scrollTop = this.$el.scrollHeight;
                }
            };
        },
        // Data
        add: function(type, msg) {
            let time = new Date().format("HH:MM");
            let obj = { type: type, msg: msg, time: time };
            this.events.push(obj);

            // Remove old events
            while (this.events.length > 10) {
                this.events.splice(0, 1);
            }
            this.window().show();
            this.window().scrollToEnd();
        }
    }
});
