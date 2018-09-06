Vue.component("search-tag", {
    props: ["name", "color", "show"],
    template: `
	<div @mouseenter="show = true;" @mouseleave="show = false;" :style="{ 'background-color': color }" class="tag">
		<div>{{name}}</div>
		<i v-show="show" @click="deleteTagFromSearch(name);" class="material-icons" style="cursor: pointer; font-size: 20px;">close</i>
	</div>`,
    methods: {
        deleteTagFromSearch: function(name) {
            this.$parent.deleteTagFromSearch(name);
        }
    }
});

Vue.component("file-tag", {
    props: ["name", "color"],
    template: `
	<div :style="{ 'background-color': color }" class="tag">
		<div>{{name}}</div>
	</div>`
});

Vue.component("tags-input", {
    props: ["name", "color"],
    methods: {
        startDrag: function(ev) {
            // Sometimes there's a bug, when user drag text, not div, so we need to check nodeName
            // If nodeName == "#text", user dragged text. We still can drop tag, but there's some graphic artifacts
            if (ev.target.nodeName == "DIV") {
                ev.dataTransfer.setData(
                    "tagName",
                    ev.target.children[0].textContent
                );
            } else if (ev.target.nodeName == "#text") {
                ev.dataTransfer.setData("tagName", ev.target.data);
            } else {
                console.error("Error: can't get the name of a tag");
            }
        }
    },
    template: `
	<div :style="{ 'background-color': color }" class="tag vertically" style="margin-bottom: 5px; margin-top: 5px;" draggable="true" @dragstart="startDrag">
		<div>{{name}}</div>
	</div>`
});

Vue.component("files", {
    props: ["file"],
    data: function() {
        return {
            hover: false
        };
    },
    methods: {
        showContextMenu: function(event, fileData) {
            this.$parent.showContextMenu(event, fileData);
        }
    },
    template: `
	<tr
		:style="[hover ? {'background-color': 'rgba(0, 0, 0, 0.1)'} : {'background-color': 'white'} ]"
		@mouseover="hover = true;"
		@mouseleave="hover = false;"
		@click.right.prevent="showContextMenu(event, file);"
		:title="file.description"
	>
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
				<file-tag
					v-for="tag in file.tags"
					:name="tag.name"
					:color="tag.color">
				</file-tag>
			</div>
		</td>
		<td>{{(file.size / (1024 * 1024)).toFixed(1)}}</td>
		<td>{{file.addTime}}</td>
	</tr>`
});
