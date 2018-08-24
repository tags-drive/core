Vue.component("tag", {
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
