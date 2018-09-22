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

// GlobalStore contains global data
var GlobalStore = {
    data: {
        allFiles: [],
        allTags: {}
    },
    updateFiles: function() {
        fetch("/api/files", {
            method: "GET",
            credentials: "same-origin"
        })
            .then(data => data.json())
            .then(files => this.setFiles(files))
            .catch(err => console.error(err));
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
            .catch(err => console.error(err));
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
                tags: () => modalWindow.showWindow().globalTags()
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
        },
        getSelectedFiles: function() {
            let files = [];
            for (i in this.$children) {
                if (this.$children[i].selected) {
                    files.push(this.$children[i].file);
                }
            }
            console.log(files);
        }
    },
    template: `
	<table :style="{'opacity': sharedState.mainBlockOpacity}" class="file-table" style="width:100%;">
		<tr style="position: sticky; top: 100px;">
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
                            console.error(text);
                            eventWindow.add(true, text);
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
                        eventWindow.add(log[i].isError, msg);
                    }
                })
                .catch(err => eventWindow.add(true, err));
        }
    }
});

// Context menu (right click on a file)
var contextMenu = new Vue({
    el: "#context-menu",
    mixins: [VueClickaway.mixin], // from vue-clickaway
    data: {
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
                    modalWindow.showWindow().renaming(this.file);
                },
                changeTags: () => {
                    this.show = false;
                    modalWindow.showWindow().fileTags(this.file);
                },
                changeDescription: () => {
                    this.show = false;
                    modalWindow.showWindow().description(this.file);
                },
                deleteFile: () => {
                    this.show = false;
                    modalWindow.showWindow().deleting(this.file);
                }
            };
        },
        // Options of context menu (select mode)
        selectMode: function() {
            return {
                changeTags: () => {
                    console.log("changeTags");
                },
                downloadFiles: () => {
                    console.log("downloadFiles");
                },
                deleteFiles: () => {
                    console.log("deleteFiles");
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
        error: "",
        sharedData: GlobalStore.data,
        // Modes
        renameMode: false,
        tagsMode: false,
        descriptionMode: false,
        deleteMode: false,
        globalTagsMode: false,
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
                renaming: file => {
                    GlobalState.showDropLayer = false;

                    this.file = file;
                    this.renameMode = true;
                    this.fileNewData.newFilename = file.filename;

                    this.show = true;
                },
                fileTags: file => {
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
                    this.tagsMode = true;

                    this.show = true;
                },
                description: file => {
                    GlobalState.showDropLayer = false;

                    this.file = file;
                    this.fileNewData.unusedTags;
                    this.descriptionMode = true;

                    this.show = true;
                },
                deleting: file => {
                    GlobalState.showDropLayer = false;

                    this.file = file;
                    this.deleteMode = true;

                    this.show = true;
                },
                globalTags: () => {
                    GlobalState.showDropLayer = false;

                    this.globalTagsMode = true;

                    this.show = true;
                }
            };
        },
        hideWindow: function() {
            this.renameMode = false;
            this.tagsMode = false;
            this.descriptionMode = false;
            this.deleteMode = false;
            this.show = false;
            this.globalTagsMode = false;

            this.error = "";

            GlobalState.showDropLayer = true;
        },
        // Drag and drop
        tagsDragAndDrop: function() {
            return {
                addToFile: ev => {
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
                delFromFile: ev => {
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
                                    console.error(text);
                                    this.error = text;
                                });
                                return;
                            }
                            // Refresh list of files
                            topBar.search().usual();
                            this.hideWindow();
                        })
                        .catch(err => {
                            this.error = err;
                            console.error(err);
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
                                    console.error(text);
                                    this.error = text;
                                });
                                return;
                            }
                            // Refresh list of files
                            topBar.search().usual();
                            this.hideWindow();
                        })
                        .catch(err => {
                            this.error = err;
                            console.error(err);
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
                                    console.error(text);
                                    this.error = text;
                                });
                                return;
                            }
                            // Refresh list of files
                            topBar.search().usual();
                            this.hideWindow();
                        })
                        .catch(err => {
                            this.error = err;
                            console.error(err);
                        });
                },
                delete: () => {
                    let params = new URLSearchParams();
                    params.append("file", this.file.filename);

                    fetch("/api/files?" + params, {
                        method: "DELETE",
                        credentials: "same-origin"
                    })
                        .then(resp => {
                            if (isErrorStatusCode(resp.status)) {
                                resp.text().then(text => {
                                    console.error(text);
                                    this.error = text;
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
                                eventWindow.add(log[i].isError, msg);
                            }
                        })
                        .catch(err => eventWindow.add(true, err));
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
                                    console.error(text);
                                    this.error = text;
                                });
                                return;
                            }

                            this.tagsAPI().delNewTag();
                            GlobalStore.updateTags();
                        })
                        .catch(err => {
                            console.error(err);
                            this.error = err;
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
                                    console.error(text);
                                    this.error = text;
                                });
                                return;
                            }

                            GlobalStore.updateTags();
                        })
                        .catch(err => {
                            console.error(err);
                            this.error = err;
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
                                    console.error(text);
                                    this.error = text;
                                });
                                return;
                            }

                            GlobalStore.updateTags();
                            // Need to update files to remove deleted tag
                            topBar.search().usual();
                            return resp.text();
                        })
                        .catch(err => {
                            console.error(err);
                            this.error = err;
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

var eventWindow = new Vue({
    el: "#events-window",
    data: {
        // States
        show: false,
        isMouseInside: false, // if isMouseInside, hideAfter isn't changed
        hideAfter: 1000 * 2, // time in milliseconds
        // UI
        opacity: 1,
        lastScrollHeight: 0,
        // Data
        /* events - array of objects:
           {
             isError: boolean,
             msg: string,
             time: string
           }
        */
        events: []
    },
    created: function() {
        const msTimeout = 50;
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
        }, msTimeout); // 50 is good enough. When 100, FPS too low.
    },
    methods: {
        // UI
        window: function() {
            return {
                show: () => {
                    this.opacity = 1;
                    this.hideAfter = 2 * 1000; // 2s
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
        add: function(isError, msg) {
            let time = new Date().format("HH:MM");
            let obj = { isError: isError, msg: msg, time: time };
            console.log(obj); // We should log obj, because there's rotation of messages
            this.events.push(obj);

            if (this.events.length > 5) {
                this.events.splice(0, 1); // remove first message
            }
            this.window().show();
            this.window().scrollToEnd();
        }
    }
});
