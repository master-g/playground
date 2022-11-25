use winit::{
  event::*,
  event_loop::{ControlFlow, EventLoop},
  window::{Window, WindowBuilder},
};

struct State {
  surface: wgpu::Surface,
  device: wgpu::Device,
  queue: wgpu::Queue,
  config: wgpu::SurfaceConfiguration,
  size: winit::dpi::PhysicalSize<u32>,
  clear_color: wgpu::Color,
}

impl State {
  // 某些 wgpu 类型需要使用异步代码才能创建
  async fn new(window: &Window) -> Self {
    let size = window.inner_size();

    // instance is the handle to the GPU
    // Backends::all means use Vulkan + Metal + DX12 + Browser WebGPU
    let instance = wgpu::Instance::new(wgpu::Backends::all());

    let surface = unsafe { instance.create_surface(window) };

    // adapter 是指向实际显卡的一个 handle， 通过它可以获取显卡的信息（例如名称，所适配到的后端等）
    // 稍后我们会使用它来创建 device 和 queue
    let adapter = instance
      .request_adapter(&wgpu::RequestAdapterOptions {
        // 性能选项， `LowPower` 和 `HighPerformance` 会影响显卡的选择, 例如在笔记本上，`LowPower` 会选择集成显卡, `HighPerformance` 会选择独立显卡
        // 当不存在 `HighPerformance` 时，会选择 `LowPower`
        power_preference: wgpu::PowerPreference::default(),
        // 要求 wgpu 所找到的适配器应当与传入的 surface 兼容
        compatible_surface: Some(&surface),
        // 强制 wgpu 选择一个能在所有硬件上工作的适配器
        force_fallback_adapter: false,
      })
      .await
      .unwrap();

    /*
      // 我们传递给 request_adapter 的选项未必能对所有设备生效，但应当能在大多数设备上可用
      // 如果 wgpu 找不到一个合适的适配器， request_adapter 会返回 None
      // 我们可以使用  `enumerate_adapters` 来获取所有可用的适配器
    let adapter = instance
      .enumerate_adapters(wgpu::Backends::all())
      .filter(|adapter| surface.get_preferred_format(adapter).is_some()) // 选择支持 surface 的适配器
      .next()
      .unwrap();
    */

    // device 和 queue
    let (device, queue) = adapter
      .request_device(
        &wgpu::DeviceDescriptor {
          // 要用的特性
          features: wgpu::Features::empty(),
          // 资源类型限制
          limits: wgpu::Limits::default(),
          label: None,
        },
        None, // Trace path (API 调用路径)
      )
      .await
      .unwrap();

    let config = wgpu::SurfaceConfiguration {
      usage: wgpu::TextureUsages::RENDER_ATTACHMENT,
      format: *surface.get_supported_formats(&adapter).first().unwrap(),
      width: size.width,
      height: size.height,
      alpha_mode: wgpu::CompositeAlphaMode::Auto,
      // 为了避免屏幕闪烁，我们使用 VSync, 也是移动设备上的最理想模式
      present_mode: wgpu::PresentMode::Fifo,
    };
    surface.configure(&device, &config);

    Self {
      surface,
      device,
      queue,
      config,
      size,
      clear_color: wgpu::Color {
        r: 0.1,
        g: 0.2,
        b: 0.3,
        a: 1.0,
      },
    }
  }

  fn resize(&mut self, new_size: winit::dpi::PhysicalSize<u32>) {
    if new_size.width > 0 && new_size.height > 0 {
      self.size = new_size;
      self.config.width = new_size.width;
      self.config.height = new_size.height;
      self.surface.configure(&self.device, &self.config);
    }
  }

  fn input(&mut self, _event: &WindowEvent) -> bool {
    match _event {
      WindowEvent::CursorMoved {
        device_id: _,
        position,
        modifiers: _,
      } => {
        let vx = position.x as f64 / self.size.width as f64;
        let vy = position.y as f64 / self.size.height as f64;
        let y = 1.0;
        let x = (y / vy) * vx;
        let z = (y / vy) * (1.0 - vx - vy);

        let mut rgb = vec![
          x * 1.656492 - y * 0.354851 - z * 0.255038,
          -x * 0.707196 + y * 1.655397 + z * 0.036152,
          x * 0.051713 - y * 0.121364 + z * 1.011530,
        ];
        rgb.iter_mut().for_each(|v| {
          // reverse gamma correction
          if *v < 0.0031308 {
            *v = 12.92 * *v;
          } else {
            *v = 1.055 * v.powf(1.0 / 2.4) - 0.055;
          }

          // bring all negative values to zero
          *v = v.max(0.0);
        });

        let mut max = f64::MIN;
        rgb.iter().for_each(|v| {
          if *v > max {
            max = *v;
          }
        });

        if max > 1.0 {
          rgb.iter_mut().for_each(|v| *v /= max);
        }

        self.clear_color = wgpu::Color {
          r: rgb[0],
          g: rgb[1],
          b: rgb[2],
          a: 1.0,
        };
      }
      _ => {}
    }
    false
  }

  fn update(&mut self) {}

  fn render(&mut self) -> Result<(), wgpu::SurfaceError> {
    let output = self.surface.get_current_texture()?;
    let view = output
      .texture
      .create_view(&wgpu::TextureViewDescriptor::default());
    let mut encoder =
      self
        .device
        .create_command_encoder(&wgpu::CommandEncoderDescriptor {
          label: Some("Render Encoder"),
        });
    {
      let _render_pass =
        encoder.begin_render_pass(&wgpu::RenderPassDescriptor {
          label: Some("Render Pass"),
          color_attachments: &[Some(wgpu::RenderPassColorAttachment {
            view: &view,
            resolve_target: None,
            ops: wgpu::Operations {
              load: wgpu::LoadOp::Clear(self.clear_color),
              store: true,
            },
          })],
          depth_stencil_attachment: None,
        });
    }

    // submit 方法能传入任何实现了 IntoIter 的参数
    self.queue.submit(std::iter::once(encoder.finish()));
    output.present();

    Ok(())
  }
}

pub fn entry() {
  // initialize global logger with env logger
  env_logger::init();

  // create window context
  let event_loop = EventLoop::new();

  // build window object
  let window = WindowBuilder::new().build(&event_loop).unwrap();

  // create our state
  let mut state = pollster::block_on(State::new(&window));

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
      Event::RedrawRequested(window_id) if window_id == window.id() => {
        state.update();
        match state.render() {
          Ok(_) => {}
          // 上下文丢失，则重新创建 surface
          Err(wgpu::SurfaceError::Lost) => state.resize(state.size),
          // 系统内存不足 (Out of Memory)，则退出
          Err(wgpu::SurfaceError::OutOfMemory) => {
            *control_flow = ControlFlow::Exit
          }
          // 其他错误，都应该在下一帧解决
          Err(e) => eprintln!("{:?}", e),
        }
      }
      Event::MainEventsCleared => {
        // RedrawRequested will only trigger once, unless we manually
        // request it.
        window.request_redraw();
      }
      Event::WindowEvent {
        ref event,
        window_id,
      } if window_id == window.id() => {
        if !state.input(event) {
          match event {
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
            WindowEvent::Resized(physical_size) => {
              state.resize(*physical_size);
            }
            WindowEvent::ScaleFactorChanged { new_inner_size, .. } => {
              // new_inner_size is &mut so winit can reuse the memory
              state.resize(**new_inner_size);
            }
            _ => {}
          }
        }
      }
      _ => {}
    }
  });
}
