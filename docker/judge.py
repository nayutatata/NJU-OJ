import os
import sys
import subprocess

def compile_code(code_file):
    compile_command = ["g++", "-o", "executable", code_file]
    compile_process = subprocess.run(compile_command, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    if compile_process.returncode != 0:
        print("Compile Error:")
        print(compile_process.stderr.decode())
        sys.exit(1)

def run_test(input_file, output_file, timeout=3):
    run_command = ["./executable"]
    with open(input_file, 'r') as infile:
        try:
            run_process = subprocess.run(run_command, stdin=infile, stdout=subprocess.PIPE, stderr=subprocess.PIPE,timeout=timeout)
        except subprocess.TimeoutExpired:
            print("TimeLimit Exceeded.")
            sys.exit(1)
    if run_process.returncode != 0:
        print("RunTime Error:")
        print(run_process.stderr.decode())
        sys.exit(1)
    output = run_process.stdout.decode()
    with open(output_file, 'r') as outfile:
        expected_output = outfile.read().strip()
    return output.strip() == expected_output

def main():
    if len(sys.argv) != 4:
        print("Usage: python judge.py code.cpp inputs outputs")
        sys.exit(1)
    
    code_file = sys.argv[1]
    inputs_folder = sys.argv[2]
    outputs_folder = sys.argv[3]
    
    # 编译代码
    compile_code(code_file)
    
    input_files = os.listdir(inputs_folder)
    output_files = os.listdir(outputs_folder)
    total_tests = min(len(input_files), len(output_files))
    passed_tests = 0
    
    for input_file, output_file in zip(input_files, output_files):
        input_path = os.path.join(inputs_folder, input_file)
        output_path = os.path.join(outputs_folder, output_file)
        if run_test(input_path, output_path):
            passed_tests += 1
    if (passed_tests==total_tests):
        print("AC")
    else:
        print("Result: "+str(passed_tests)+"/"+str(total_tests))

if __name__ == "__main__":
    main()
