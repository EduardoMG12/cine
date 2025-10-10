module.exports = {
	"*.{js,jsx,ts,tsx}": (filenames) =>
		filenames.map((filename) => `biome format --write "${filename}"`),

	"api_v2/**/*.go": (filenames) => {
		const commands = [];
		// Format each file individually with correct path
		filenames.forEach((filename) => {
			commands.push(`gofmt -w "${filename}"`);
		});
		// Skip go vet for now due to module resolution issues
		// TODO: Fix go vet after module setup is complete
		return commands;
	},

	"flutter_app/**/*.dart": (filenames) => {
		const commands = [];
		// Format each file individually with correct path
		filenames.forEach((filename) => {
			commands.push(`flutter format "${filename}"`);
		});
		return commands;
	},
};
