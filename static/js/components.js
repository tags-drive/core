// Tags in Search bar
Vue.component("search-tag", {
    props: {
        tag: Object
    },
    data: function() {
        return {
            show: false
        };
    },
    template: `
	<div @mouseenter="show = true;" @mouseleave="show = false;" :style="{ 'background-color': tag.color }" class="tag">
		<div>{{tag.name}}</div>
		<i v-show="show" @click="deleteTagFromSearch();" class="material-icons" style="cursor: pointer; font-size: 20px;">close</i>
	</div>`,
    methods: {
        deleteTagFromSearch: function() {
            this.$parent.input().tags.delete(this.tag.id);
        }
    }
});

Vue.component("suggestion-tag", {
    props: {
        tag: Object
    },
    methods: {
        add: function() {
            this.$parent.input().tags.add(this.tag.id);
        }
    },
    template: `
	<div class="top-bar__search__suggestion-tag" @click="add">
		<div :style="{ 'background-color': tag.color }" class="tag" style="margin: 0; cursor: pointer;" :title="'id - ' + tag.id">
			{{tag.name}}
		</div>
	</div>`
});

// For drag and drop input
Vue.component("tags-input", {
    props: ["tag"],
    methods: {
        startDrag: function(ev) {
            ev.dataTransfer.setData("tagName", this.tag.id);
        }
    },
    template: `
	<div :style="{ 'background-color': tag.color }" class="tag vertically" style="margin-bottom: 5px; margin-top: 5px;" draggable="true" @dragstart="startDrag">
		<div>{{tag.name}}</div>
	</div>`
});

const validTagName = /^[\w\d- ]{1,20}$/;
const validColor = /^#[\dabcdef]{6}$/;

// For tags editing
Vue.component("modifying-tags", {
    props: {
        tag: Object,
        isNewTag: String // only new tag
    },
    data: function() {
        return {
            newName: this.tag.name,
            newColor: this.tag.color,
            isChanged: this.isNewTag !== true ? false : true, // isNewTag wasn't passed,
            isError: false,
            isDeleted: false
        };
    },
    destroyed: function() {
        // Called, when window is closed
        // We delete a tag only after closing the window
        // It lets us to undo the file deleting
        if (this.isDeleted) {
            this.$parent.tagsAPI().del(this.tag.id);
        }
    },
    methods: {
        check: function() {
            if (this.tag.name == this.newName && this.tag.color == this.newColor && this.isNewTag !== true) {
                // Can skip, if name and color weren't changed
                this.isChanged = false;
                this.isError = false;
                return;
            }
            this.isChanged = true;

            if (this.newName.length == 0 || validTagName.exec(this.newName) === null) {
                this.isError = true;
                return;
            }
            if (validColor.exec(this.newColor) === null) {
                this.isError = true;
                return;
            }

            this.isError = false;
        },
        generateRandomColor: function() {
            if (this.isDeleted) {
                return;
            }
            this.isChanged = true;
            this.isError = false; // we can't generate an invalid color
            this.newColor = "#" + Math.floor(Math.random() * 16777215).toString(16);
        },
        // API
        save: function() {
            if (this.isError || !this.isChanged) {
                return;
            }

            if (this.isNewTag) {
                // Need to create, not to change
                this.$parent.tagsAPI().add(this.newName, this.newColor);
            } else {
                this.$parent.tagsAPI().change(this.tag.id, this.newName, this.newColor);
            }

            this.isChanged = false;
        },
        del: function() {
            if (this.isNewTag) {
                // Delete tag right now
                this.$parent.tagsAPI().delNewTag();
                return;
            }

            this.isDeleted = true;
        },
        recover: function() {
            this.isDeleted = false;
        }
    },
    template: `
	<div style="display: inline-flex; margin-bottom: 5px; width: 95%;">
		<div style="width: 2px; height: 20px; margin-right: 3px;" class="vertically">
			<div v-if="isDeleted" style="height: 20px; border-left: 2px solid white;"></div>
			<div v-else-if="isError"	style="height: 20px; border-left: 2px solid red;"></div>
			<div v-else-if="isChanged" style="height: 20px; border-left: 2px solid blue;"></div>
		</div>
		
		<div style="width: 35%; display: flex;">
			<div :style="{ 'background-color': newColor }" class="tag">
				<div>{{newName}}</div>
			</div>
		</div>

		<input @input="check" type="text" maxlength="20" :disabled="isDeleted" v-model="newName" style="width: 35%; margin-right: 10px;">

		<input @input="check" type="text" :disabled="isDeleted" v-model="newColor" style="width: 15%; margin-right: 5px;">

		<i class="material-icons btn"
			title="Generate a new color"
			@click="generateRandomColor"
			style="margin-right: 10px;"
			:style="[isDeleted ? {'opacity': '0.3', 'background-color': 'white', 'cursor': 'default'} : {'opacity': '1'}]">cached</i>

		<div style="display: flex;">
			<i class="material-icons btn" title="Save" @click="save" style="margin-right: 5px;" 
			:style="[isError || isDeleted || !this.isChanged ? {'opacity': '0.3', 'background-color': 'white', 'cursor': 'default'} : {'opacity': '1'}]">done</i>

			<i v-if="!isDeleted"
				class="material-icons btn"
				title="Delete"
				@click="del"
				:style="[isDeleted ? {'opacity': '0.3', 'background-color': 'white', 'cursor': 'default'} : {'opacity': '1'}]"
			>delete</i>
			<i v-else
				class="material-icons btn"
				title="Undo"
				@click="recover"
			>undo</i>
		</div>
	</div>`
});

// For selected mod
Vue.component("selected-add-tag", {
    props: {
        tag: Object
    },
    data: function() {
        return {
            shouldAdd: false
        };
    },
    destroyed: function() {
        if (this.shouldAdd) {
            this.$parent.filesAPI().addSelectedFilesTag(this.tag.id);
        }
    },
    template: `
	<div style="display: flex; margin-right: auto; margin-left: auto; margin-bottom: 5px; position: relative;">
		<div style="width: 200px; display: flex">
			<div :style="{ 'background-color': tag.color }" class="tag" style="margin: 0;">
				<div>{{tag.name}}</div>
			</div>
		</div>
		<div style="position: absolute; right: 0;">
			<input v-model="shouldAdd" type="checkbox" style="width: 20px; height: 20px; right: 0;" title="Add tag">
		</div>
	</div>`
});

Vue.component("selected-delete-tag", {
    props: {
        tag: Object
    },
    data: function() {
        return {
            shouldDelete: false
        };
    },
    destroyed: function() {
        if (this.shouldDelete) {
            this.$parent.filesAPI().deleteSelectedFilesTag(this.tag.id);
        }
    },
    template: `
	<div style="display: flex; margin-right: auto; margin-left: auto; margin-bottom: 5px; position: relative;">
		<div style="width: 200px; display: flex">
			<div :style="{ 'background-color': tag.color }" class="tag" style="margin: 0;">
				<div>{{tag.name}}</div>
			</div>
		</div>
		<div style="position: absolute; right: 0;">
			<input v-model="shouldDelete" type="checkbox" style="width: 20px; height: 20px; right: 0;" title="Add tag">
		</div>
	</div>`
});

// Files in Main block
Vue.component("files", {
    props: {
        file: Object,
        allTags: Object
    },
    data: function() {
        return {
            hover: false,
            selected: false
        };
    },
    methods: {
        showContextMenu: function(event) {
            contextMenu.setFile(this.file);
            contextMenu.showMenu(event.x, event.y);
        },
        toggleSelect: function() {
            // We can skip changing this.selected, because a checkbox is bound to this.selected

            // The function is called after changing this.selected
            if (this.selected) {
                this.$parent.selectFile();
            } else {
                this.$parent.unselectFile();
            }
        },
        /* For the parent */
        select: function() {
            this.selected = true;
        },
        unselect: function() {
            this.selected = false;
        }
    },
    template: `
	<tr
		:style="[hover || selected ? {'background-color': 'rgba(0, 0, 0, 0.1)'} : {'background-color': 'white'} ]"
		@mouseover="hover = true;"
		@mouseleave="hover = false;"
		@click.right.prevent="showContextMenu($event, file);"
		:title="file.description"
	>
		<td style="text-align: center; width: 30px;">
			<input type="checkbox" @change="toggleSelect" v-model="selected" style="height: 15px; width: 15px;">
		</td>
		<td v-if="file.type == 'image'" style="width: 30px;">
			<img :src="file.preview" style="width: 30px;">
		</td>
		<td v-else style="width: 30px; text-align: center;">
			<img :src="'/ext/' + file.filename.split('.').pop()" style="width: 30px;">
		</td>	
		<td style="width: 200px;">
			<div class="filename" :title="file.filename">
				{{file.filename}}
			</div>
		</td>
		<td>
			<div style="display: flex;">
				<div v-for="id in file.tags" :style="{ 'background-color': allTags[id].color }" class="tag">
					<div>{{allTags[id].name}}</div>
				</div>
			</div>
		</td>
		<td>{{(file.size / (1024 * 1024)).toFixed(1)}}</td>
		<td>{{file.addTime}}</td>
	</tr>`
});
