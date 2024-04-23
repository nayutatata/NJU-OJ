package judger

/*
	pipeline:
	1. 拿到答案字符串
	2. 拿到题目的输入输出字符串
	3. 将答案、输入输出写入container的文件中
	4. 在container中执行python脚本，得到答案
	5. 返回python脚本的stdout的内容
*/
import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"unicode"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

const (
	image_name string = "real_final"
)

func Judge_samples(answer string, inputs, outputs []string) string {
	dockercli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}
	hconfig := container.HostConfig{}
	hconfig.CPUPeriod = 100000
	hconfig.CPUQuota = 100000
	hconfig.Memory = 256 * 1024 * 1024
	hconfig.MemorySwappiness = new(int64)
	/*
		hconfig.Mounts = []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: "D:\\important\\study\\网络应用开发\\NJU-OJ\\docker\\t",
				Target: "/app",
			},
		}*/
	hconfig.ReadonlyRootfs = false
	ctx := context.Background()
	resp, err := dockercli.ContainerCreate(ctx, &container.Config{
		Image: image_name,
		//Entrypoint: []string{"echo", "hello"},
		User: "root",
	}, &hconfig, nil, nil, "")
	if err != nil {
		log.Fatal(err)
	}
	defer dockercli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{
		Force: true,
	})
	err = dockercli.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		log.Fatal(err)
	}
	// 将answer拷贝到文件中
	//debug(dockercli, ctx)
	codepath := "./code.cpp"
	err = create_file(codepath, answer, dockercli, ctx, resp.ID)
	//err = dockercli.CopyToContainer(ctx, resp.ID, codepath, strings.NewReader("1"), types.CopyToContainerOptions{})

	//err = create_file(codepath, answer, dockercli, ctx, resp.ID)
	if err != nil {
		dockercli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{
			Force: true,
		})
		log.Fatal(err)
	}
	// 将inputs拷贝到文件夹中
	dir := "./inputs"
	for i, input := range inputs {
		input_path := dir + fmt.Sprintf("/%d", i)
		/*
			err = dockercli.CopyToContainer(ctx, resp.ID, input_path, strings.NewReader(input), types.CopyToContainerOptions{
				AllowOverwriteDirWithFile: true,
			})*/
		err = create_file(input_path, input, dockercli, ctx, resp.ID)
		if err != nil {
			dockercli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{
				Force: true,
			})
			log.Fatal(err)
		}
	}
	// 将outputs拷贝到文件夹中
	dir = "./outputs"
	for i, output := range outputs {
		output_path := dir + fmt.Sprintf("/%d", i)
		/*
			err = dockercli.CopyToContainer(ctx, resp.ID, output_path, strings.NewReader(output), types.CopyToContainerOptions{
				AllowOverwriteDirWithFile: true,
			})
		*/
		err = create_file(output_path, output, dockercli, ctx, resp.ID)
		if err != nil {
			dockercli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{
				Force: true,
			})
			log.Fatal(err)
		}
	}
	execConfig := types.ExecConfig{
		Cmd:          []string{"python3", "./judge.py", "./code.cpp", "./inputs", "./outputs"},
		AttachStdout: true,
		AttachStderr: true,
	}
	execresp, err := dockercli.ContainerExecCreate(ctx, resp.ID, execConfig)
	if err != nil {
		log.Fatal(err)
	}

	result, err := dockercli.ContainerExecAttach(ctx, execresp.ID, types.ExecStartCheck{})
	if err != nil {
		log.Fatal(err)
	}
	defer result.Close()

	var judge_res strings.Builder
	_, err = io.Copy(&judge_res, result.Reader)
	if err != nil {
		log.Fatal(err)
	}

	return regular(judge_res.String())
}
func create_file(path string, content string, cli *client.Client, ctx context.Context, containerID string) error {
	command := fmt.Sprintf("echo \"%s\" > %s", content, path)
	//command := "touch code.cpp"
	execConfig := types.ExecConfig{
		User: "root",
		//Cmd:  []string{"echo", fmt.Sprintf("\"%s\"", content), ">", path},
		Cmd: []string{"sh", "-c", command},
		Tty: true,
	}
	execresp, err := cli.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return err
	}
	err = cli.ContainerExecStart(ctx, execresp.ID, types.ExecStartCheck{})
	return err
}
func debug(cli *client.Client, ctx context.Context) {
	containers, err := cli.ContainerList(ctx, container.ListOptions{
		All: false, // 只列出正在运行的容器
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, con := range containers {
		fmt.Print(con.Names)
	}
}
func regular(s string) string {
	return strings.TrimLeftFunc(s, func(r rune) bool {
		return unicode.IsControl(r) || unicode.IsSpace(r) || !unicode.IsPrint(r)
	})
}
