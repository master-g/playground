mod editor;
mod terminal;

pub use editor::Position;
pub use terminal::Terminal;

use editor::Editor;

fn main() {
  Editor::default().run();
}
