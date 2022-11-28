use wgpu::{include_wgsl, util::DeviceExt};
use winit::{
  event::*,
  event_loop::{ControlFlow, EventLoop},
  window::{Window, WindowBuilder},
};

#[repr(C)]
#[derive(Copy, Clone, Debug, bytemuck::Pod, bytemuck::Zeroable)]
struct Vertex {
  position: [f32; 3],
  color: [f32; 3],
}

const VERTICES: &[Vertex] = &[
  Vertex {
    position: [-0.0868241, 0.49240386, 0.0],
    color: [0.5, 0.0, 0.5],
  }, // A
  Vertex {
    position: [-0.49513406, 0.06958647, 0.0],
    color: [0.5, 0.0, 0.5],
  }, // B
  Vertex {
    position: [-0.21918549, -0.44939706, 0.0],
    color: [0.5, 0.0, 0.5],
  }, // C
  Vertex {
    position: [0.35966998, -0.3473291, 0.0],
    color: [0.5, 0.0, 0.5],
  }, // D
  Vertex {
    position: [0.44147372, 0.2347359, 0.0],
    color: [0.5, 0.0, 0.5],
  }, // E
];

const INDICES: &[u16] = &[0, 1, 4, 1, 2, 4, 2, 3, 4];

impl Vertex {
  const ATTRIBS: [wgpu::VertexAttribute; 2] = wgpu::vertex_attr_array![
    0 => Float32x3,
    1 => Float32x3,
  ];

  fn desc<'a>() -> wgpu::VertexBufferLayout<'a> {
    wgpu::VertexBufferLayout {
      array_stride: std::mem::size_of::<Vertex>() as wgpu::BufferAddress,
      step_mode: wgpu::VertexStepMode::Vertex,
      attributes: &Self::ATTRIBS,
    }
  }
}

struct State {
  surface: wgpu::Surface,
  device: wgpu::Device,
  queue: wgpu::Queue,
  config: wgpu::SurfaceConfiguration,

  render_pipeline: wgpu::RenderPipeline,

  vertex_buffers: Vec<wgpu::Buffer>,
  index_buffers: Vec<wgpu::Buffer>,
  num_indices: Vec<u32>,

  clear_color: wgpu::Color,
  challenge_mode: bool,

  size: winit::dpi::PhysicalSize<u32>,
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
          limits: if cfg!(target_arch = "wasm32") {
            wgpu::Limits::downlevel_webgl2_defaults()
          } else {
            wgpu::Limits::default()
          },
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

    // pipeline
    let render_pipeline_layout =
      device.create_pipeline_layout(&wgpu::PipelineLayoutDescriptor {
        label: Some("Render Pipeline Layout"),
        bind_group_layouts: &[],
        push_constant_ranges: &[],
      });

    let shader = device
      .create_shader_module(include_wgsl!("../../static/shader/chpt04.wgsl"));

    let render_pipeline =
      device.create_render_pipeline(&wgpu::RenderPipelineDescriptor {
        label: Some("Render Pipeline"),
        layout: Some(&render_pipeline_layout),
        vertex: wgpu::VertexState {
          module: &shader,
          entry_point: "vs_main",
          buffers: &[Vertex::desc()],
        },
        primitive: wgpu::PrimitiveState {
          topology: wgpu::PrimitiveTopology::TriangleList,
          strip_index_format: None,
          front_face: wgpu::FrontFace::Ccw,
          cull_mode: Some(wgpu::Face::Back),
          // 如果将该字段设置为除了 Fill 之外的任何值，都需要 Features::NON_FILL_POLYGON_MODE
          polygon_mode: wgpu::PolygonMode::Fill,
          // 需要 Features::DEPTH_CLIP_ENABLE
          unclipped_depth: false,
          // 需要 Features::CONSERVATIVE_RASTERIZATION
          conservative: false,
        },
        depth_stencil: None,
        multisample: wgpu::MultisampleState {
          count: 1,
          mask: !0,
          alpha_to_coverage_enabled: false,
        },
        fragment: Some(wgpu::FragmentState {
          module: &shader,
          entry_point: "fs_main",
          targets: &[Some(wgpu::ColorTargetState {
            format: config.format,
            blend: Some(wgpu::BlendState::REPLACE),
            write_mask: wgpu::ColorWrites::ALL,
          })],
        }),
        multiview: None,
      });

    // buffer
    let vertex_buffer =
      device.create_buffer_init(&wgpu::util::BufferInitDescriptor {
        label: Some("Vertex Buffer"),
        contents: bytemuck::cast_slice(VERTICES),
        usage: wgpu::BufferUsages::VERTEX,
      });

    let index_buffer =
      device.create_buffer_init(&wgpu::util::BufferInitDescriptor {
        label: Some("Index Buffer"),
        contents: bytemuck::cast_slice(INDICES),
        usage: wgpu::BufferUsages::INDEX,
      });

    let num_vertices = 16;
    let angle = std::f32::consts::PI * 2.0 / num_vertices as f32;
    let challenge_verts = (0..num_vertices)
      .map(|i| {
        let theta = angle * i as f32;
        Vertex {
          position: [0.5 * theta.cos(), -0.5 * theta.sin(), 0.0],
          color: [(1.0 + theta.cos()) / 2.0, (1.0 + theta.sin()) / 2.0, 1.0],
        }
      })
      .collect::<Vec<_>>();

    let num_triangles = num_vertices - 2;
    let challenge_indices = (1u16..num_triangles + 1)
      .into_iter()
      .flat_map(|i| vec![i + 1, i, 0])
      .collect::<Vec<_>>();
    let num_challenge_indices = challenge_indices.len() as u32;

    let challenge_vertex_buffer =
      device.create_buffer_init(&wgpu::util::BufferInitDescriptor {
        label: Some("Challenge Vertex Buffer"),
        contents: bytemuck::cast_slice(&challenge_verts),
        usage: wgpu::BufferUsages::VERTEX,
      });
    let challenge_index_buffer =
      device.create_buffer_init(&wgpu::util::BufferInitDescriptor {
        label: Some("Challenge Index Buffer"),
        contents: bytemuck::cast_slice(&challenge_indices),
        usage: wgpu::BufferUsages::INDEX,
      });

    Self {
      surface,
      device,
      queue,
      config,

      render_pipeline,

      vertex_buffers: vec![vertex_buffer, challenge_vertex_buffer],
      index_buffers: vec![index_buffer, challenge_index_buffer],
      num_indices: vec![INDICES.len() as u32, num_challenge_indices],

      clear_color: wgpu::Color {
        r: 0.1,
        g: 0.2,
        b: 0.3,
        a: 1.0,
      },
      challenge_mode: false,

      size,
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

  #[allow(deprecated)]
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
      WindowEvent::KeyboardInput {
        device_id: _,
        input,
        is_synthetic: _,
      } => {
        if input.state == ElementState::Pressed {
          if let Some(keycode) = input.virtual_keycode {
            match keycode {
              VirtualKeyCode::Space => {
                self.challenge_mode = !self.challenge_mode;
              }
              _ => {}
            }
          }
        }
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
      let mut render_pass =
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

      render_pass.set_pipeline(&self.render_pipeline);
      let i = if self.challenge_mode { 1 } else { 0 };
      render_pass.set_vertex_buffer(0, self.vertex_buffers[i].slice(..));
      render_pass.set_index_buffer(
        self.index_buffers[i].slice(..),
        wgpu::IndexFormat::Uint16,
      );
      render_pass.draw_indexed(0..self.num_indices[i], 0, 0..1);
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
