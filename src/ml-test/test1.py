import tensorflow as tf
import matplotlib.pyplot as plt

# 实际的线
TRUE_W = 3.0
TRUE_B = 2.0

NUM_EXAMPLES = 1000

# 随机向量x
x = tf.random.normal(shape=[NUM_EXAMPLES])

# 生成噪声
noise = tf.random.normal(shape=[NUM_EXAMPLES])

# 计算y
y = x * TRUE_W + TRUE_B + noise

class MyModel(tf.Module):
  def __init__(self, **kwargs):
    super().__init__(**kwargs)
    # 初始化权重值为`5.0`，偏差值为`0.0`
    # 实际项目中，应该随机初始化
    self.w = tf.Variable(5.0)
    self.b = tf.Variable(0.0)

  def __call__(self, x):
    return self.w * x + self.b

# 计算整个批次的单个损失值
def loss(target_y, predicted_y):
  return tf.reduce_mean(tf.square(target_y - predicted_y))

# 给定一个可调用的模型，输入，输出和学习率...
def train(model, x, y, learning_rate):

  with tf.GradientTape() as t:
    # 可训练变量由GradientTape自动跟踪
    current_loss = loss(y, model(x))

  # 使用GradientTape计算相对于W和b的梯度
  dw, db = t.gradient(current_loss, [model.w, model.b])

  # 减去由学习率缩放的梯度
  model.w.assign_sub(learning_rate * dw)
  model.b.assign_sub(learning_rate * db)

# 定义用于训练的循环
def training_loop(model, x, y):

  for epoch in epochs:
    # 用单个大批次处理更新模型
    train(model, x, y, learning_rate=0.1)

    # 在更新之前进行跟踪
    Ws.append(model.w.numpy())
    bs.append(model.b.numpy())
    current_loss = loss(y, model(x))

    print("Epoch %2d: W=%1.2f b=%1.2f, loss=%2.5f" %
          (epoch, Ws[-1], bs[-1], current_loss))

model = MyModel()

# 收集W值和b值的历史记录以供以后绘制
Ws, bs = [], []
epochs = range(10)

print("Starting: W=%1.2f b=%1.2f, loss=%2.5f" %
      (model.w, model.b, loss(y, model(x))))

# 开始训练
training_loop(model, x, y)

print("Ending: W=%1.2f b=%1.2f, loss=%2.5f" %
      (model.w, model.b, loss(y, model(x))))

class MyModelKeras(tf.keras.Model):
  def __init__(self, **kwargs):
    super().__init__(**kwargs)
    # 初始化权重为`5.0`，偏差为`0.0`
    # 实际中应该随机初始化该值
    self.w = tf.Variable(5.0)
    self.b = tf.Variable(0.0)

  def __call__(self, x, **kwargs):
    return self.w * x + self.b

keras_model = MyModelKeras()

# 编译设置培训参数
keras_model.compile(
    # 默认情况下，fit()调用tf.function()。
    # Debug时你可以关闭这一功能，但是现在是打开的。
    run_eagerly=False,

    # 使用内置的优化器，配置为对象
    optimizer=tf.keras.optimizers.SGD(learning_rate=0.1),

    # Keras内置MSE
    # 您也可以使用损失函数像上面一样进行定义
    loss=tf.keras.losses.mean_squared_error,
)

keras_model.fit(x, y, epochs=10, batch_size=1000)

print("Ending: W=%1.2f b=%1.2f, loss=%2.5f" %
      (keras_model.w, keras_model.b, loss(y, keras_model(x))))