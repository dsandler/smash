extern crate smash;
use smash::term::Term;
use smash::view;
use std::cell::RefCell;
use std::rc::Rc;

fn main() {
    view::init();
    let win = view::Win::new();

    {
        let mut win = win.borrow_mut();
        let font_extents = {
            let ctx = win.create_cairo();
            Term::get_font_metrics(&ctx)
        };

        win.resize(80 * font_extents.max_x_advance as i32,
                   25 * font_extents.height as i32);

        let term = Term::new(win.dirty_cb.clone(),
                             font_extents,
                             &["bash"],
                             Box::new(|| {}));

        win.child = Rc::new(RefCell::new(term));
        win.show();
    }

    view::main();
}