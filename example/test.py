import docker
from docker.types import Mount


client = docker.from_env()
image_name = 'registry.dp.tech/dptech/abacus:3.1.0'
client.images.pull(image_name)
abacus_path = '/root/test_rg/abacus_input/000_I3Pb1C1N2H5/lcao'
container_path = '/root/'

mount = Mount(target=container_path, source=abacus_path,
              type='bind', read_only=False)


container = client.containers.create(
    image=image_name,
    command='sh',
    tty=True,
    stdin_open=True,
    name='my_container',
    mounts=[mount]
)

container.start()
command = "sh -c 'cd /root && mpirun -n 8 abacus'"
exec_instance = container.exec_run(cmd=command, stdout=True, stderr=True)
output = exec_instance.output.decode('utf-8')
print(f"Output of '{command}':\n{output}")

container.stop()
container.remove()

client.close()
