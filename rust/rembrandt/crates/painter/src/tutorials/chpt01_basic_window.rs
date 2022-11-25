use winit::{
  event::*,
  event_loop::{ControlFlow, EventLoop},
  window::WindowBuilder,
};

pub fn entry() {
  // initialize global logger with env logger
  env_logger::init();

  // create window context
  let event_loop = EventLoop::new();

  // build window object
  let window = WindowBuilder::new().build(&event_loop).unwrap();

  // start event loop
  event_loop.run(move |event, _, control_flow| {
    // ControlFlow::Poll continuously runs the event loop, even if the OS hasn't
    // dispatched any events. This is ideal for games and similar applications.
    // control_flow.set_poll();

    // ControlFlow::Wait pauses the event loop if no events are available to process.
    // This is ideal for non-game applications that only update in response to user
    // input, and uses significantly less power/CPU time than ControlFlow::Poll.
    control_flow.set_wait();

    match event {
      // handle event
      Event::WindowEvent {
        ref event,
        window_id,
      } if window_id == window.id() => match event {
        WindowEvent::CloseRequested
        | WindowEvent::KeyboardInput {
          input:
            KeyboardInput {
              state: ElementState::Pressed,
              virtual_keycode: Some(VirtualKeyCode::Escape),
              ..
            },
          ..
        } => *control_flow = ControlFlow::Exit,
        _ => {}
      },
      _ => {}
    }
  });
}
