const path = require("path");

module.exports = {
	"*.{js,jsx,ts,tsx}": (filenames) =>
		filenames.map((filename) => `biome format --write "${filename}"`),

	"api_v2/**/*.go": (filenames) => {
		const relativeFilepaths = filenames
			.map((filename) => path.relative("api_v2", filename))
			.join(" ");
		return [`gofmt -w ${relativeFilepaths}`, "go vet ./..."];
	},

	"flutter_app/**/*.dart": (filenames) => {
		const relativeFilepaths = filenames
			.map((filename) => path.relative("flutter_app", filename))
			.join(" ");
		return [`flutter format ${relativeFilepaths}`];
	},
};
